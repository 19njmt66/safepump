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

    // Helper to complete the incubation phase (>= 20% progress = 160M tokens sold)
    function _completeIncubation() internal {
        // Under 5 ETH curve:
        // A buy of 0.015 ETH is under the 2% limit (buys ~14M tokens)
        // 30 buys of 0.015 ETH = 0.45 ETH, which is enough to reach 20% progress (requires 0.2 ETH)
        for (uint256 i = 0; i < 30; i++) {
            address user = address(uint160(200 + i));
            vm.deal(user, 10 ether);
            vm.roll(block.number + 1);
            vm.prank(user);
            factory.buy{value: 0.015 ether}(tokenAddress);
        }
        assertTrue(token.incubationComplete());
    }

    function test_createToken() public {
        vm.prank(creator);
        // Default deployment without initial buy
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        assertEq(token.name(), "SafePump Coin");
        assertEq(token.symbol(), "SPC");
        assertEq(token.creator(), creator);
        assertEq(token.factory(), address(factory));
        
        // Creator gets 0 tokens
        assertEq(token.balanceOf(creator), 0);
        
        // Factory holds the entire total supply (1,000,000,000 tokens)
        assertEq(token.balanceOf(address(factory)), 1_000_000_000 * 10**18);
    }

    function test_createTokenWithBonus() public {
        vm.prank(creator);
        // Deploy with creator bonus buy (0.02127 ether)
        tokenAddress = factory.createToken{value: 0.02127 ether}("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        // Creator gets 50,000,000 tokens (5% bonus allocation)
        assertEq(token.balanceOf(creator), 50_000_000 * 10**18);

        // Factory holds 950,000,000 tokens
        assertEq(token.balanceOf(address(factory)), 950_000_000 * 10**18);

        // Verify state
        (,,,, bool migrated, bool hasBonus) = factory.tokens(tokenAddress);
        assertFalse(migrated);
        assertTrue(hasBonus);
    }

    function test_launchBlockCooldown() public {
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");

        // First buy in current block succeeds (using small amount to stay under 2% max wallet)
        vm.prank(buyer1);
        factory.buy{value: 0.015 ether}(tokenAddress);

        // Second buy in the same block reverts
        vm.prank(buyer2);
        vm.expectRevert("SafePump: One buy per block limit active");
        factory.buy{value: 0.015 ether}(tokenAddress);

        // Move to next block
        vm.roll(block.number + 1);

        // Next block buy succeeds
        vm.prank(buyer2);
        factory.buy{value: 0.015 ether}(tokenAddress);
    }

    function test_maxWalletLimit() public {
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        // 0.015 ETH buy gets ~14M tokens (under 2% = 20M limit)
        vm.prank(buyer1);
        factory.buy{value: 0.015 ether}(tokenAddress);
        assertTrue(token.balanceOf(buyer1) > 0);

        // Another buy by buyer1 in the next block that crosses the 2% limit
        vm.roll(block.number + 1);
        vm.prank(buyer1);
        vm.expectRevert("SafePump: Exceeds max wallet limit (2%)");
        factory.buy{value: 0.015 ether}(tokenAddress); // 14M + 13.8M = 27.8M (exceeds 20M)
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
        assertTrue(token.balanceOf(buyer1) > 20_000_000 * 10**18); // holds > 2% supply
    }

    function test_sellLimit() public {
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        // Complete incubation so max wallet limit is disabled
        _completeIncubation();

        // Roll block beyond launch block cooldown
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

    function test_creatorBonusAndLP() public {
        // Creator deploys with 0.02127 ether to activate the 5% bonus
        vm.prank(creator);
        tokenAddress = factory.createToken{value: 0.02127 ether}("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        uint256 creatorBalanceBeforeMigration = creator.balance;

        // Complete incubation
        _completeIncubation();

        // Buy remaining tokens to trigger migration
        // Curve requires 5 ETH total.
        // We've already raised: 0.02127 ETH (creator) + 30 * 0.015 ETH (incubation) = 0.47127 ETH
        // Remaining needed: 5.00 - 0.47127 = 4.52873 ETH
        // Buy with 4.6 ETH to easily trigger completion (excess will be refunded)
        vm.deal(buyer1, 10 ether);
        vm.roll(block.number + 1);
        vm.prank(buyer1);
        factory.buy{value: 4.6 ether}(tokenAddress);

        assertTrue(token.migrationComplete());

        // Verify creator received 0 refund (since refund is removed)
        assertEq(creator.balance - creatorBalanceBeforeMigration, 0, "Creator should not get a refund");

        // Verify Uniswap LP receives the exact mathematically calculated WETH (5 ETH - creator bonus curve impact)
        address wethAddress = router.WETH();
        address pair = MockUniswapV2Factory(uniFactory).getPair(tokenAddress, wethAddress);
        assertEq(MockWETH(wethAddress).balanceOf(pair), 4.808473125 ether, "LP ETH balance mismatch");

        // Verify Uniswap LP receives 20% of tokens (200,000,000)
        assertEq(token.balanceOf(pair), 200_000_000 * 10**18, "LP token balance mismatch");
    }

    function test_financialMechanisms() public {
        // 1. Create a token
        vm.prank(creator);
        tokenAddress = factory.createToken("SafePump Coin", "SPC");
        token = SafePumpToken(payable(tokenAddress));

        // Complete incubation phase to lift the 2% max wallet limit
        _completeIncubation();

        uint256 initTokensSold;
        uint256 initEthRaised;
        {
            (,, uint256 ts, uint256 er,, ) = factory.tokens(tokenAddress);
            initTokensSold = ts;
            initEthRaised = er;
        }

        uint256 tokensBought;
        
        // 2. Perform a buy swap from buyer1 of 0.04 ETH (keeps bought tokens under 2% sell limit)
        {
            uint256 initCreatorBalance = creator.balance;
            uint256 initFeeRecipientBalance = feeRecipient.balance;

            vm.prank(buyer1);
            factory.buy{value: 0.04 ether}(tokenAddress);

            // Verify balances after buy (1% fee = 0.0004 ether, 100% to platform feeRecipient)
            assertEq(creator.balance - initCreatorBalance, 0, "Creator should get 0 fees");
            assertEq(feeRecipient.balance - initFeeRecipientBalance, 0.0004 ether, "Platform fee recipient fee mismatch");

            (,, uint256 ts, uint256 er,, ) = factory.tokens(tokenAddress);
            assertEq(er - initEthRaised, 0.0396 ether, "ETH raised mismatch on buy");
            tokensBought = ts - initTokensSold;
        }

        // 3. Perform a sell swap from buyer1
        vm.roll(block.number + 1);

        {
            uint256 expectedEthOut = factory.getAmountOutEth(tokenAddress, tokensBought);
            uint256 sellFee = expectedEthOut / 100;
            uint256 expectedEthToUser = expectedEthOut - sellFee;

            uint256 balCreatorBeforeSell = creator.balance;
            uint256 balFeeRecipientBeforeSell = feeRecipient.balance;
            uint256 balBuyerBeforeSell = buyer1.balance;

            vm.prank(buyer1);
            token.approve(address(factory), tokensBought);
            
            vm.prank(buyer1);
            factory.sell(tokenAddress, tokensBought);

            // Verify fee distribution on sell
            assertEq(creator.balance - balCreatorBeforeSell, 0, "Creator should get 0 fees on sell");
            assertEq(feeRecipient.balance - balFeeRecipientBeforeSell, sellFee, "Platform should get 100% fees on sell");
            
            // Verify buyer received net ETH
            assertEq(buyer1.balance - balBuyerBeforeSell, expectedEthToUser, "Buyer net ETH output mismatch");
        }

        // Verify factory token state is reset to incubation state (with 1 wei tolerance for rounding)
        {
            (,, uint256 tokensSoldAfter, uint256 ethRaisedAfter,, ) = factory.tokens(tokenAddress);
            assertApproxEqAbs(tokensSoldAfter, initTokensSold, 1, "Tokens sold mismatch after sell");
            assertApproxEqAbs(ethRaisedAfter, initEthRaised, 1, "ETH raised mismatch after sell");
        }
    }
}
