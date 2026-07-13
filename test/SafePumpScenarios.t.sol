// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {SafePumpToken} from "../src/SafePumpToken.sol";
import {SafePumpFactory} from "../src/SafePumpFactory.sol";

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

contract SafePumpScenariosTest is Test {
    SafePumpFactory public factory;
    MockUniswapV2Router public router;
    MockUniswapV2Factory public uniFactory;
    MockWETH public weth;

    address public owner = address(100);
    address public feeRecipient = address(101);
    address public creator = address(102);
    address public buyer1 = address(103);

    function setUp() public {
        vm.deal(owner, 1000 ether);
        vm.deal(creator, 1000 ether);
        vm.deal(buyer1, 1000 ether);

        uniFactory = new MockUniswapV2Factory();
        weth = new MockWETH();
        router = new MockUniswapV2Router(address(uniFactory), address(weth));

        vm.prank(owner);
        factory = new SafePumpFactory(address(router), feeRecipient, owner);
    }

    function test_ScenarioA_NoBonus() public {
        console.log("----------------------------------------------------------------");
        console.log("SCENARIO A : CREATOR DOES NOT ACTIVATE BONUS (Deploy with 0 ETH)");
        console.log("----------------------------------------------------------------");

        // 1. Creator deploys token
        uint256 initCreatorEth = creator.balance;
        vm.prank(creator);
        address tokenAddress = factory.createToken("No Bonus Coin", "NBC");
        SafePumpToken token = SafePumpToken(payable(tokenAddress));

        // 2. Retail traders buy during incubation (< 20% progress = 160M tokens)
        for (uint256 i = 0; i < 22; i++) {
            address user = address(uint160(200 + i));
            vm.deal(user, 1 ether);
            vm.roll(block.number + 1);
            vm.prank(user);
            factory.buy{value: 0.01 ether}(tokenAddress);
        }
        assertTrue(token.incubationComplete(), "Incubation should be complete");

        // 3. Whale buyer1 completes bonding curve (requires 5 ETH total)
        uint256 initBuyer1Eth = buyer1.balance;
        vm.roll(block.number + 1);
        vm.prank(buyer1);
        factory.buy{value: 4.9 ether}(tokenAddress); // Excess will be refunded automatically

        assertTrue(token.migrationComplete(), "Migration should be complete");
        address pair = MockUniswapV2Factory(uniFactory).getPair(tokenAddress, address(weth));

        // 4. Print final results
        console.log("Balances des acteurs apres Migration :");
        console.log("  1. Le Createur (Dev) :");
        console.log("     - Solde ETH (wei)  :", creator.balance);
        console.log("     - Solde NBC (tokens):", token.balanceOf(creator) / 10**18);
        console.log("     - ETH depense (wei):", initCreatorEth - creator.balance);

        console.log("  2. La Plateforme (Treasure) :");
        console.log("     - Solde ETH (wei)  :", feeRecipient.balance);

        console.log("  3. Les Traders (Buyer1 - Whale) :");
        console.log("     - Solde ETH (wei)  :", buyer1.balance);
        console.log("     - Solde NBC (tokens):", token.balanceOf(buyer1) / 10**18);
        console.log("     - ETH depense (wei):", initBuyer1Eth - buyer1.balance);

        console.log("  4. Le Pool Uniswap (Pair LP) :");
        console.log("     - Solde WETH (wei) :", MockWETH(address(weth)).balanceOf(pair));
        console.log("     - Solde NBC (tokens):", token.balanceOf(pair) / 10**18);
        console.log("----------------------------------------------------------------");
    }

    function test_ScenarioB_WithBonus() public {
        console.log("----------------------------------------------------------------");
        console.log("SCENARIO B : CREATOR ACTIVATES BONUS (Deploy with 0.02127 ETH)");
        console.log("----------------------------------------------------------------");

        // 1. Creator deploys token and purchases the 2% package (0.02127 ETH)
        uint256 initCreatorEth = creator.balance;
        vm.prank(creator);
        address tokenAddress = factory.createToken{value: 0.02127 ether}("Bonus Coin", "BC");
        SafePumpToken token = SafePumpToken(payable(tokenAddress));

        // 2. Retail traders buy during incubation (< 20% progress = 160M tokens)
        for (uint256 i = 0; i < 18; i++) {
            address user = address(uint160(200 + i));
            vm.deal(user, 1 ether);
            vm.roll(block.number + 1);
            vm.prank(user);
            factory.buy{value: 0.01 ether}(tokenAddress);
        }
        assertTrue(token.incubationComplete(), "Incubation should be complete");

        // 3. Whale buyer1 completes bonding curve (requires 5 ETH total in curve)
        uint256 initBuyer1Eth = buyer1.balance;
        vm.roll(block.number + 1);
        vm.prank(buyer1);
        factory.buy{value: 4.9 ether}(tokenAddress); // Excess will be refunded automatically

        assertTrue(token.migrationComplete(), "Migration should be complete");
        address pair = MockUniswapV2Factory(uniFactory).getPair(tokenAddress, address(weth));

        // 4. Print final results
        console.log("Balances des acteurs apres Migration :");
        console.log("  1. Le Createur (Dev) :");
        console.log("     - Solde ETH (wei)  :", creator.balance);
        console.log("     - Solde BC (tokens) :", token.balanceOf(creator) / 10**18);
        console.log("     - ETH depense (wei):", initCreatorEth - creator.balance);

        console.log("  2. La Plateforme (Treasure) :");
        console.log("     - Solde ETH (wei)  :", feeRecipient.balance);

        console.log("  3. Les Traders (Buyer1 - Whale) :");
        console.log("     - Solde ETH (wei)  :", buyer1.balance);
        console.log("     - Solde BC (tokens) :", token.balanceOf(buyer1) / 10**18);
        console.log("     - ETH depense (wei):", initBuyer1Eth - buyer1.balance);

        console.log("  4. Le Pool Uniswap (Pair LP) :");
        console.log("     - Solde WETH (wei) :", MockWETH(address(weth)).balanceOf(pair));
        console.log("     - Solde BC (tokens) :", token.balanceOf(pair) / 10**18);
        console.log("----------------------------------------------------------------");
    }
}
