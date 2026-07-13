// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {SafePumpFactory} from "../src/SafePumpFactory.sol";
import {SafePumpToken} from "../src/SafePumpToken.sol";

// Mocks for Uniswap V2 Router & Factory
contract MockWETH {
    string public name = "Wrapped Ether";
    string public symbol = "WETH";
    uint8 public decimals = 18;
    mapping(address => uint256) public balanceOf;

    function deposit() public payable {
        balanceOf[msg.sender] += msg.value;
    }

    function transfer(address to, uint256 value) public returns (bool) {
        require(balanceOf[msg.sender] >= value, "MockWETH: Insufficient balance");
        balanceOf[msg.sender] -= value;
        balanceOf[to] += value;
        return true;
    }
}

contract MockUniswapV2Pair {
    address public token0;
    address public token1;

    constructor(address _token0, address _token1) {
        token0 = _token0;
        token1 = _token1;
    }
}

contract MockUniswapV2Factory {
    mapping(address => mapping(address => address)) public getPair;

    function createPair(address tokenA, address tokenB) external returns (address pair) {
        pair = address(new MockUniswapV2Pair(tokenA, tokenB));
        getPair[tokenA][tokenB] = pair;
        getPair[tokenB][tokenA] = pair;
    }
}

contract MockUniswapV2Router {
    address public factory;
    address public WETH;

    constructor(address _factory, address _weth) {
        factory = _factory;
        WETH = _weth;
    }

    function addLiquidityETH(
        address token,
        uint amountTokenDesired,
        uint amountTokenMin,
        uint amountETHMin,
        address to,
        uint deadline
    ) external payable returns (uint amountToken, uint amountETH, uint liquidity) {
        MockWETH(WETH).deposit{value: msg.value}();
        
        address pair = MockUniswapV2Factory(factory).getPair(token, WETH);
        if (pair == address(0)) {
            pair = MockUniswapV2Factory(factory).createPair(token, WETH);
        }

        SafePumpToken(payable(token)).transferFrom(msg.sender, pair, amountTokenDesired);
        MockWETH(WETH).transfer(pair, msg.value);

        return (amountTokenDesired, msg.value, 100 ether);
    }
}

contract DeploySafePumpLocal is Script {
    function run() external {
        uint256 deployerPrivateKey = 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80; // Account #0
        address deployer = vm.addr(deployerPrivateKey);

        vm.startBroadcast(deployerPrivateKey);

        // Deploy Mocks
        MockUniswapV2Factory uniFactory = new MockUniswapV2Factory();
        MockWETH weth = new MockWETH();
        MockUniswapV2Router router = new MockUniswapV2Router(address(uniFactory), address(weth));

        // Deploy SafePumpFactory pointing to local mock router
        SafePumpFactory factory = new SafePumpFactory(
            address(router),
            deployer, // Platform Fee Recipient (Account #0)
            deployer  // Owner
        );

        vm.stopBroadcast();

        console.log("MOCK_UNISWAP_V2_ROUTER:", address(router));
        console.log("LOCAL_FACTORY_ADDRESS:", address(factory));
    }
}
