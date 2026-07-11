// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {SafePumpToken} from "./SafePumpToken.sol";

interface IUniswapV2Router02 {
    function factory() external pure returns (address);
    function WETH() external pure returns (address);
    function addLiquidityETH(
        address token,
        uint amountTokenDesired,
        uint amountTokenMin,
        uint amountETHMin,
        address to,
        uint deadline
    ) external payable returns (uint amountToken, uint amountETH, uint liquidity);
}

interface IUniswapV2Factory {
    function getPair(address tokenA, address tokenB) external view returns (address pair);
}

/**
 * @title SafePumpFactory
 * @dev Factory contract to deploy SafePumpTokens and manage their bonding curve lifecycle,
 * fee collection, and automatic Uniswap V2 migration.
 */
contract SafePumpFactory is Ownable {
    // Constant Bonding Curve Parameters
    uint256 public constant INITIAL_VIRTUAL_TOKEN_RESERVES = 941_176_470 * 10**18;
    uint256 public constant INITIAL_VIRTUAL_ETH_RESERVES = 3 * 10**18;
    uint256 public constant MAX_BONDING_CURVE_TOKENS = 800_000_000 * 10**18;
    uint256 public constant TARGET_ETH_RAISED = 17 * 10**18;

    struct TokenInfo {
        address tokenAddress;
        address creator;
        uint256 tokensSold;
        uint256 ethRaised;
        bool migrated;
    }

    // Configurable addresses
    address public uniswapV2Router;
    address public feeRecipient;

    // Token tracking
    mapping(address => TokenInfo) public tokens;
    address[] public allTokens;

    // Fee tracking
    mapping(address => uint256) public pendingCreatorFees;
    uint256 public claimablePlatformFees;

    // Events
    event TokenCreated(address indexed token, address indexed creator, string name, string symbol);
    event Buy(address indexed token, address indexed buyer, uint256 tokenAmount, uint256 ethAmount, uint256 fee);
    event Sell(address indexed token, address indexed seller, uint256 tokenAmount, uint256 ethAmount, uint256 fee);
    event Migrated(address indexed token, address indexed pair, uint256 ethAmount, uint256 tokenAmount);

    constructor(
        address _uniswapV2Router,
        address _feeRecipient,
        address _initialOwner
    ) Ownable(_initialOwner) {
        require(_uniswapV2Router != address(0), "SafePump: Invalid router");
        require(_feeRecipient != address(0), "SafePump: Invalid fee recipient");
        uniswapV2Router = _uniswapV2Router;
        feeRecipient = _feeRecipient;
    }

    /**
     * @notice Deploy a new SafePumpToken
     */
    function createToken(string memory name, string memory symbol) external returns (address) {
        SafePumpToken newToken = new SafePumpToken(name, symbol, msg.sender);
        address tokenAddress = address(newToken);

        // Transfer 5% creator vesting allocation to the creator
        newToken.transfer(msg.sender, 50_000_000 * 10**18);

        tokens[tokenAddress] = TokenInfo({
            tokenAddress: tokenAddress,
            creator: msg.sender,
            tokensSold: 0,
            ethRaised: 0,
            migrated: false
        });

        allTokens.push(tokenAddress);

        emit TokenCreated(tokenAddress, msg.sender, name, symbol);
        return tokenAddress;
    }

    /**
     * @notice Buy tokens on the bonding curve using ETH
     */
    function buy(address tokenAddress) external payable {
        TokenInfo storage info = tokens[tokenAddress];
        require(info.tokenAddress != address(0), "SafePump: Token not found");
        require(!info.migrated, "SafePump: Token already migrated");
        require(msg.value > 0, "SafePump: ETH must be greater than zero");

        uint256 ethAmount = msg.value;
        uint256 remainingTokens = MAX_BONDING_CURVE_TOKENS - info.tokensSold;

        uint256 fee;
        uint256 netEthAmount;
        uint256 tokensToBuy;
        uint256 excessEth = 0;

        // Calculate temporary values
        uint256 tempFee = ethAmount / 100;
        uint256 tempNetEth = ethAmount - tempFee;
        uint256 tempTokensToBuy = getAmountOutTokens(tokenAddress, tempNetEth);

        if (tempTokensToBuy >= remainingTokens) {
            // Bonding curve completes
            tokensToBuy = remainingTokens;

            uint256 currentTokenReserves = INITIAL_VIRTUAL_TOKEN_RESERVES - info.tokensSold;
            uint256 currentEthReserves = INITIAL_VIRTUAL_ETH_RESERVES + info.ethRaised;
            uint256 targetTokenReserves = currentTokenReserves - remainingTokens;

            uint256 k = currentTokenReserves * currentEthReserves;
            uint256 netEthNeeded = (k / targetTokenReserves) - currentEthReserves;

            // Total ETH needed = netEthNeeded * 100 / 99
            uint256 totalEthNeeded = (netEthNeeded * 100) / 99;
            require(ethAmount >= totalEthNeeded, "SafePump: Insufficient ETH to complete");

            fee = totalEthNeeded - netEthNeeded;
            netEthAmount = netEthNeeded;
            excessEth = ethAmount - totalEthNeeded;
        } else {
            tokensToBuy = tempTokensToBuy;
            fee = tempFee;
            netEthAmount = tempNetEth;
        }

        // Update state
        info.tokensSold += tokensToBuy;
        info.ethRaised += netEthAmount;

        // Distribute fees
        uint256 creatorShare = (fee * 40) / 100;
        uint256 platformShare = fee - creatorShare;
        _distributeFees(info.creator, creatorShare, platformShare);

        // Transfer tokens to the buyer
        SafePumpToken(payable(tokenAddress)).transfer(msg.sender, tokensToBuy);

        emit Buy(tokenAddress, msg.sender, tokensToBuy, netEthAmount, fee);

        // Check and set incubation completion if progress >= 20% (160M tokens)
        SafePumpToken token = SafePumpToken(payable(tokenAddress));
        if (!token.incubationComplete() && info.tokensSold >= 160_000_000 * 10**18) {
            token.setIncubationComplete();
        }

        // Refund excess ETH
        if (excessEth > 0) {
            (bool success, ) = msg.sender.call{value: excessEth}("");
            require(success, "SafePump: Refund failed");
        }

        // Migrate if bonding curve target reached
        if (info.tokensSold >= MAX_BONDING_CURVE_TOKENS) {
            _migrate(tokenAddress);
        }
    }

    /**
     * @notice Sell tokens back to the bonding curve for ETH
     */
    function sell(address tokenAddress, uint256 tokenAmount) external {
        TokenInfo storage info = tokens[tokenAddress];
        require(info.tokenAddress != address(0), "SafePump: Token not found");
        require(!info.migrated, "SafePump: Token already migrated");
        require(tokenAmount > 0, "SafePump: Token amount must be greater than zero");
        require(tokenAmount <= info.tokensSold, "SafePump: Cannot sell more than total sold");

        uint256 ethOutNet = getAmountOutEth(tokenAddress, tokenAmount);

        // Calculate fee (1%)
        uint256 fee = ethOutNet / 100;
        uint256 ethToUser = ethOutNet - fee;

        // Update state
        info.tokensSold -= tokenAmount;
        info.ethRaised -= ethOutNet;

        // Distribute fees
        uint256 creatorShare = (fee * 40) / 100;
        uint256 platformShare = fee - creatorShare;
        _distributeFees(info.creator, creatorShare, platformShare);

        // Transfer tokens from user to Factory
        SafePumpToken(payable(tokenAddress)).transferFrom(msg.sender, address(this), tokenAmount);

        // Send ETH to user
        (bool success, ) = msg.sender.call{value: ethToUser}("");
        require(success, "SafePump: ETH transfer failed");

        emit Sell(tokenAddress, msg.sender, tokenAmount, ethOutNet, fee);
    }

    /**
     * @notice Calculate tokens received for a given ETH input
     */
    function getAmountOutTokens(address tokenAddress, uint256 ethIn) public view returns (uint256) {
        TokenInfo memory info = tokens[tokenAddress];
        require(info.tokenAddress != address(0), "SafePump: Token not found");
        require(!info.migrated, "SafePump: Token already migrated");

        uint256 currentTokenReserves = INITIAL_VIRTUAL_TOKEN_RESERVES - info.tokensSold;
        uint256 currentEthReserves = INITIAL_VIRTUAL_ETH_RESERVES + info.ethRaised;

        uint256 numerator = currentTokenReserves * ethIn;
        uint256 denominator = currentEthReserves + ethIn;
        return numerator / denominator;
    }

    /**
     * @notice Calculate ETH received for a given token input
     */
    function getAmountOutEth(address tokenAddress, uint256 tokensIn) public view returns (uint256) {
        TokenInfo memory info = tokens[tokenAddress];
        require(info.tokenAddress != address(0), "SafePump: Token not found");
        require(!info.migrated, "SafePump: Token already migrated");

        uint256 currentTokenReserves = INITIAL_VIRTUAL_TOKEN_RESERVES - info.tokensSold;
        uint256 currentEthReserves = INITIAL_VIRTUAL_ETH_RESERVES + info.ethRaised;

        uint256 numerator = currentEthReserves * tokensIn;
        uint256 denominator = currentTokenReserves + tokensIn;
        return numerator / denominator;
    }

    /**
     * @notice Internal migration to Uniswap V2 and burning LP tokens
     */
    function _migrate(address tokenAddress) internal {
        TokenInfo storage info = tokens[tokenAddress];
        info.migrated = true;

        SafePumpToken token = SafePumpToken(payable(tokenAddress));

        // Fees at migration
        uint256 platformMigrationFee = 0.03 ether;
        uint256 creatorBonus = 0.02 ether;
        uint256 totalMigrationFee = platformMigrationFee + creatorBonus;

        require(info.ethRaised >= totalMigrationFee, "SafePump: Insufficient ETH raised");

        uint256 ethForLP = info.ethRaised - totalMigrationFee;
        uint256 tokenAmountForLP = 150_000_000 * 10**18; // 15% LP allocation

        // Distribute migration fees
        (bool successPlatform, ) = feeRecipient.call{value: platformMigrationFee}("");
        if (!successPlatform) {
            claimablePlatformFees += platformMigrationFee;
        }

        (bool successCreator, ) = info.creator.call{value: creatorBonus}("");
        if (!successCreator) {
            pendingCreatorFees[info.creator] += creatorBonus;
        }

        // Approve Uniswap Router to spend LP tokens
        token.approve(uniswapV2Router, tokenAmountForLP);

        // Add liquidity on Uniswap V2 Router
        // Mint LP tokens directly to address(0) to burn them permanently (anti-LP-rug)
        IUniswapV2Router02(uniswapV2Router).addLiquidityETH{value: ethForLP}(
            tokenAddress,
            tokenAmountForLP,
            tokenAmountForLP, // amountTokenMin
            ethForLP,         // amountETHMin
            address(0),       // LP tokens recipient (0x0 = Burned)
            block.timestamp
        );

        address factoryAddress = IUniswapV2Router02(uniswapV2Router).factory();
        address wethAddress = IUniswapV2Router02(uniswapV2Router).WETH();
        address pair = IUniswapV2Factory(factoryAddress).getPair(tokenAddress, wethAddress);

        // Mark complete on token contract
        token.setMigrationComplete(pair, uniswapV2Router);

        emit Migrated(tokenAddress, pair, ethForLP, tokenAmountForLP);
    }

    /**
     * @notice Internal fee distribution helper
     */
    function _distributeFees(address creator, uint256 creatorShare, uint256 platformShare) internal {
        if (platformShare > 0) {
            (bool success, ) = feeRecipient.call{value: platformShare}("");
            if (!success) {
                claimablePlatformFees += platformShare;
            }
        }

        if (creatorShare > 0) {
            (bool success, ) = creator.call{value: creatorShare}("");
            if (!success) {
                pendingCreatorFees[creator] += creatorShare;
            }
        }
    }

    /**
     * @notice Claim pending fees for creators
     */
    function claimCreatorFees() external {
        uint256 amount = pendingCreatorFees[msg.sender];
        require(amount > 0, "SafePump: No fees to claim");
        pendingCreatorFees[msg.sender] = 0;
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "SafePump: Claim failed");
    }

    /**
     * @notice Claim pending fees for platform
     */
    function claimPlatformFees() external onlyOwner {
        uint256 amount = claimablePlatformFees;
        require(amount > 0, "SafePump: No fees to claim");
        claimablePlatformFees = 0;
        (bool success, ) = feeRecipient.call{value: amount}("");
        require(success, "SafePump: Claim failed");
    }

    // Owner administration
    function setFeeRecipient(address _feeRecipient) external onlyOwner {
        require(_feeRecipient != address(0), "SafePump: Invalid address");
        feeRecipient = _feeRecipient;
    }

    function setUniswapV2Router(address _router) external onlyOwner {
        require(_router != address(0), "SafePump: Invalid address");
        uniswapV2Router = _router;
    }

    // Return list of all tokens
    function getAllTokens() external view returns (address[] memory) {
        return allTokens;
    }

    receive() external payable {}
}
