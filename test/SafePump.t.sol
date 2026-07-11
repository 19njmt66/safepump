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

contract SafePumpTest is Test {
    SafePumpFactory public factory;
    MockUniswapV2Router public router;
    MockUniswapV2Factory public uniFactory;
    MockWETH public weth;

    address public owner = address(100);
    address public feeRecipient = address(101);
    address public creator = address(102);
    address public buyer1 = address(103);
    address public buyer2 = address(104);

    SafePumpToken public token;
    address public tokenAddress;

    function setUp() public {
        vm.deal(owner, 1000 ether);
        vm.deal(creator, 1000 ether);
        vm.deal(buyer1, 1000 ether);
        vm.deal(buyer2, 1000 ether);

        // Deploy mocks
        uniFactory = new MockUniswapV2Factory();
        weth = new MockWETH();
        router = new MockUniswapV2Router(address(uniFactory), address(weth));

        // Deploy Factory
        vm.prank(owner);
        factory = new SafePumpFactory(address(router), feeRecipient, owner);
    }

    // Helper to complete the incubation phase (>= 20% progress)
    function _completeIncubation() internal {
        for (uint256 i = 0; i < 30; i++) {
            address user = address(uint160(200 + i));
            vm.deal(user, 10 ether);
            vm.roll(block.number + 1);
            vm.prank(user);
            factory.buy{value: 0.03 ether}(tokenAddress);
        }
        assertTrue(token.incubationComplete());
    }

    function test_createToken() public {
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        assertEq(token.name(), "SafePump Coin");
        assertEq(token.symbol(), "SPC");
        assertEq(token.creator(), creator);
        assertEq(token.factory(), address(factory));
        
        // Creator gets 50,000,000 tokens (5%)
        assertEq(token.balanceOf(creator), 50_000_000 * 10**18);
        
        // Factory holds the rest of total supply (950,000,000 tokens)
        assertEq(token.balanceOf(address(factory)), 950_000_000 * 10**18);
    }

    function test_launchBlockCooldown() public {
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");

        // First buy in current block succeeds (using small amount to stay under 1.5% max wallet)
        vm.prank(buyer1);
        factory.buy{value: 0.03 ether}(tokenAddress);

        // Second buy in the same block reverts
        vm.prank(buyer2);
        vm.expectRevert("SafePump: One buy per block limit active");
        factory.buy{value: 0.03 ether}(tokenAddress);

        // Move to next block
        vm.roll(block.number + 1);

        // Next block buy succeeds
        vm.prank(buyer2);
        factory.buy{value: 0.03 ether}(tokenAddress);
    }

    function test_maxWalletLimit() public {
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        // 0.03 ETH buy:
        // netEth = 0.03 * 0.99 = 0.0297 ETH
        // tokens = 941,176,470 * 0.0297 / 3.0297 = 9.22M tokens (under 15M limit)
        vm.prank(buyer1);
        factory.buy{value: 0.03 ether}(tokenAddress);
        assertTrue(token.balanceOf(buyer1) > 0);

        // Another buy by buyer1 that crosses the 15M limit
        vm.roll(block.number + 1);
        vm.prank(buyer1);
        vm.expectRevert("SafePump: Exceeds max wallet limit (1.5%)");
        factory.buy{value: 0.03 ether}(tokenAddress); // 9.2M + 9.2M = 18.4M (exceeds 15M)
    }

    function test_incubationCompletion() public {
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        _completeIncubation();

        // Now, max wallet limit is disabled. A single user can buy a large amount without revert.
        vm.roll(block.number + 1);
        vm.prank(buyer1);
        factory.buy{value: 1 ether}(tokenAddress);
        assertTrue(token.balanceOf(buyer1) > 15_000_000 * 10**18); // holds > 1.5% supply
    }

    function test_sellLimit() public {
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        // Complete incubation so max wallet limit is disabled
        _completeIncubation();

        // Roll block beyond launch block cooldown (launch block + 10)
        vm.roll(block.number + 10);

        // Simulate providing tokens to buyer1 from Factory (Factory is excluded from sell limits)
        vm.prank(address(factory));
        token.transfer(buyer1, 25_000_000 * 10**18);

        // Buyer1 tries to transfer 21,000,000 tokens (exceeds 2% 24h sell limit)
        vm.prank(buyer1);
        vm.expectRevert("SafePump: Exceeds 24h sell limit (2%)");
        token.transfer(buyer2, 21_000_000 * 10**18);

        // Buyer1 transfers 10,000,000 tokens (succeeds)
        vm.prank(buyer1);
        token.transfer(buyer2, 10_000_000 * 10**18);

        // Transferring another 11,000,000 within the same 24h period fails (total 21M > 20M)
        vm.prank(buyer1);
        vm.expectRevert("SafePump: Exceeds 24h sell limit (2%)");
        token.transfer(buyer2, 11_000_000 * 10**18);

        // Warp time by 24h + 1s
        skip(1 days + 1);

        // Transfer succeeds now
        vm.prank(buyer1);
        token.transfer(buyer2, 10_000_000 * 10**18);
    }

    function test_creatorVesting() public {
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        // 1. Pre-migration: Creator cannot transfer any vested tokens (reverts since lockedAmount = 50M)
        vm.prank(creator);
        vm.expectRevert("SafePump: Vesting locked");
        token.transfer(buyer1, 1);

        // 2. Complete incubation phase so max wallet limit is disabled for buyer1
        _completeIncubation();

        // 3. buyer1 completes bonding curve to migrate (since max wallet is disabled, they can buy everything)
        vm.deal(buyer1, 20 ether);
        vm.roll(block.number + 1);
        vm.prank(buyer1);
        factory.buy{value: 18 ether}(tokenAddress);

        assertTrue(token.migrationComplete());
        uint256 migrationTime = token.migrationTime();

        // Creator's balance is exactly their 50M vested allocation!
        assertEq(token.balanceOf(creator), 50_000_000 * 10**18);

        // 4. Post-migration: Creator still cannot transfer immediately (cooldown 48h active)
        vm.prank(creator);
        vm.expectRevert("SafePump: Vesting locked");
        token.transfer(buyer2, 1);

        // Warp to 48 hours post-migration
        vm.warp(migrationTime + 48 hours);

        // Day 0: 0 tokens unlocked. Still reverts
        vm.prank(creator);
        vm.expectRevert("SafePump: Vesting locked");
        token.transfer(buyer2, 1);

        // Warp to 48 hours + 1 day
        vm.warp(migrationTime + 48 hours + 1 days);

        // Day 1: 1,000,000 tokens unlocked. Let's transfer 500,000 (succeeds)
        vm.prank(creator);
        token.transfer(buyer2, 500_000 * 10**18);
        assertEq(token.balanceOf(buyer2), 500_000 * 10**18);

        // Transferring another 600,000 reverts (total 1.1M > 1M unlocked)
        vm.prank(creator);
        vm.expectRevert("SafePump: Vesting locked");
        token.transfer(buyer2, 600_000 * 10**18);

        // Warp to 48 hours + 50 days (full unlock)
        vm.warp(migrationTime + 48 hours + 50 days);

        // Creator transfers 15,000,000 tokens (succeeds, under 2% = 20M limit)
        vm.prank(creator);
        token.transfer(buyer2, 15_000_000 * 10**18);

        // Creator transfers another 15,000,000 in the same day (fails with sell limit)
        vm.prank(creator);
        vm.expectRevert("SafePump: Exceeds 24h sell limit (2%)");
        token.transfer(buyer2, 15_000_000 * 10**18);

        // Skip 1 day
        skip(1 days + 1);

        // Creator transfers another 15,000,000 (succeeds)
        vm.prank(creator);
        token.transfer(buyer2, 15_000_000 * 10**18);

        // Skip 1 day
        skip(1 days + 1);

        // Creator transfers another 15,000,000 (succeeds)
        vm.prank(creator);
        token.transfer(buyer2, 15_000_000 * 10**18);

        // Skip 1 day
        skip(1 days + 1);

        // Creator transfers the remaining 4,500,000 (succeeds)
        uint256 remaining = token.balanceOf(creator);
        vm.prank(creator);
        token.transfer(buyer2, remaining);
        
        assertEq(token.balanceOf(creator), 0);
    }

    function test_financialMechanisms() public {
        // 1. Create a token
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        // Complete incubation phase to lift the 1.5% max wallet limit
        _completeIncubation();

        uint256 initTokensSold;
        uint256 initEthRaised;
        {
            (,, uint256 ts, uint256 er, ) = factory.tokens(tokenAddress);
            initTokensSold = ts;
            initEthRaised = er;
        }

        uint256 tokensBought;
        
        // 2. Perform a buy swap from buyer1 of 0.05 ETH
        {
            uint256 initCreatorBalance = creator.balance;
            uint256 initFeeRecipientBalance = feeRecipient.balance;

            vm.prank(buyer1);
            factory.buy{value: 0.05 ether}(tokenAddress);

            // Verify balances after buy (1% fee = 0.0005 ether, 40% creator = 0.0002, 60% platform = 0.0003)
            assertEq(creator.balance - initCreatorBalance, 0.0002 ether, "Creator share on buy mismatch");
            assertEq(feeRecipient.balance - initFeeRecipientBalance, 0.0003 ether, "Platform share on buy mismatch");

            (,, uint256 ts, uint256 er, ) = factory.tokens(tokenAddress);
            assertEq(er - initEthRaised, 0.0495 ether, "ETH raised mismatch on buy");
            tokensBought = ts - initTokensSold;
        }

        // 3. Perform a sell swap from buyer1
        // To sell, we first need to roll the block since we bought in this block (cooldown)
        vm.roll(block.number + 1);

        {
            // Expected gross ETH output calculation
            uint256 expectedEthOut = factory.getAmountOutEth(tokenAddress, tokensBought);
            // Fee on sell is 1% of expectedEthOut: expectedEthOut / 100
            uint256 sellFee = expectedEthOut / 100;
            uint256 expectedEthToUser = expectedEthOut - sellFee;
            
            uint256 creatorShareSell = (sellFee * 40) / 100;
            uint256 platformShareSell = sellFee - creatorShareSell;

            // Save balances before sell
            uint256 balCreatorBeforeSell = creator.balance;
            uint256 balFeeRecipientBeforeSell = feeRecipient.balance;
            uint256 balBuyerBeforeSell = buyer1.balance;

            // Approve and sell
            vm.prank(buyer1);
            token.approve(address(factory), tokensBought);
            
            vm.prank(buyer1);
            factory.sell(tokenAddress, tokensBought);

            // Verify fee distribution on sell
            assertEq(creator.balance - balCreatorBeforeSell, creatorShareSell, "Creator share on sell mismatch");
            assertEq(feeRecipient.balance - balFeeRecipientBeforeSell, platformShareSell, "Platform share on sell mismatch");
            
            // Verify buyer received net ETH
            assertEq(buyer1.balance - balBuyerBeforeSell, expectedEthToUser, "Buyer net ETH output mismatch");
        }

        // Verify factory token state is reset to incubation state (with 1 wei tolerance for rounding)
        {
            (,, uint256 tokensSoldAfter, uint256 ethRaisedAfter, ) = factory.tokens(tokenAddress);
            assertApproxEqAbs(tokensSoldAfter, initTokensSold, 1, "Tokens sold mismatch after sell");
            assertApproxEqAbs(ethRaisedAfter, initEthRaised, 1, "ETH raised mismatch after sell");
        }
    }
}
