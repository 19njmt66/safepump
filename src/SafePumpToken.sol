// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/**
 * @title SafePumpToken
 * @dev Custom ERC-20 token for the SafePump Launchpad with embedded security rules:
 * - Block cooldown during launch (first 10 blocks: max 1 purchase transaction per block globally).
 * - Max wallet holding limit of 1.5% during incubation phase (deactivates at 20% progress).
 * - 24-hour rolling sell limit of 2% of total supply (anti-panic dump).
 * - Creator vesting (5% supply locked, unlocks linearly at 0.1% per day starting 48h after migration).
 */
contract SafePumpToken is ERC20 {
    // Constant Supply Metrics
    uint256 public constant TOTAL_SUPPLY = 1_000_000_000 * 10**18; // 1 Billion tokens
    uint256 public constant MAX_WALLET_LIMIT = 15_000_000 * 10**18; // 1.5% of total supply
    uint256 public constant SELL_LIMIT_24H = 20_000_000 * 10**18; // 2% of total supply
    uint256 public constant CREATOR_ALLOCATION = 50_000_000 * 10**18; // 5% of total supply
    uint256 public constant VESTING_RELEASE_RATE_PER_DAY = 1_000_000 * 10**18; // 0.1% of total supply per day

    // Core addresses
    address public immutable factory;
    address public immutable creator;

    // Launch & Incubation State
    uint256 public immutable launchBlock;
    uint256 public lastBuyBlock;
    bool public incubationComplete;
    bool public migrationComplete;

    // Post-migration details
    uint256 public migrationTime;
    address public uniswapV2Pair;
    address public uniswapV2Router;

    // Exclusions
    mapping(address => bool) public isExcludedFromMaxWallet;
    mapping(address => bool) public isExcludedFromSellLimit;

    // Sell limit records
    struct SellRecord {
        uint128 amountSold;
        uint128 lastWindowStart;
    }
    mapping(address => SellRecord) public sellRecords;

    // Events
    event IncubationCompleted();
    event MigrationCompleted(address indexed pair, address indexed router);

    modifier onlyFactory() {
        require(msg.sender == factory, "SafePump: Only factory can call");
        _;
    }

    constructor(
        string memory name,
        string memory symbol,
        address _creator
    ) ERC20(name, symbol) {
        require(_creator != address(0), "SafePump: Invalid creator address");
        
        factory = msg.sender;
        creator = _creator;
        launchBlock = block.number;

        // Setup exclusions
        isExcludedFromMaxWallet[factory] = true;
        isExcludedFromMaxWallet[creator] = true;
        isExcludedFromMaxWallet[address(this)] = true;
        isExcludedFromMaxWallet[address(0)] = true;
        isExcludedFromMaxWallet[address(0xdead)] = true;

        isExcludedFromSellLimit[factory] = true;
        isExcludedFromSellLimit[address(this)] = true;
        isExcludedFromSellLimit[address(0)] = true;
        isExcludedFromSellLimit[address(0xdead)] = true;

        // Mint total supply to the factory (which handles distribution)
        _mint(factory, TOTAL_SUPPLY);
    }

    /**
     * @notice Signal that incubation phase has completed (progression >= 20%)
     * @dev Called only by the Factory to disable the 1.5% max wallet limit
     */
    function setIncubationComplete() external onlyFactory {
        require(!incubationComplete, "SafePump: Already complete");
        incubationComplete = true;
        emit IncubationCompleted();
    }

    /**
     * @notice Signal that migration to Uniswap V2 has completed
     * @dev Disables max wallet limits permanently and registers Uniswap Pair/Router
     */
    function setMigrationComplete(address _pair, address _router) external onlyFactory {
        require(!migrationComplete, "SafePump: Already migrated");
        require(_pair != address(0) && _router != address(0), "SafePump: Invalid addresses");

        migrationComplete = true;
        migrationTime = block.timestamp;
        uniswapV2Pair = _pair;
        uniswapV2Router = _router;

        // Exclude exchange components
        isExcludedFromMaxWallet[_pair] = true;
        isExcludedFromMaxWallet[_router] = true;
        isExcludedFromSellLimit[_pair] = true;
        isExcludedFromSellLimit[_router] = true;

        emit MigrationCompleted(_pair, _router);
    }

    /**
     * @dev Overridden transfer logic to enforce the security rules of SafePump
     */
    function _update(
        address from,
        address to,
        uint256 value
    ) internal override {
        // 1. Vesting check if the sender is the creator
        if (from == creator) {
            uint256 lockedAmount = getLockedCreatorAmount();
            require(balanceOf(creator) >= value + lockedAmount, "SafePump: Vesting locked");
        }

        // 2. Cooldown check: 1 buy transaction per block in the first 10 blocks (from deployment)
        if (from == factory && !isExcludedFromMaxWallet[to]) {
            if (block.number < launchBlock + 10) {
                require(lastBuyBlock != block.number, "SafePump: One buy per block limit active");
                lastBuyBlock = block.number;
            }
        }

        // 3. Max Wallet Limit check (active during incubation, which is < 20% progress, and pre-migration)
        if (!incubationComplete && !migrationComplete) {
            if (!isExcludedFromMaxWallet[to]) {
                require(balanceOf(to) + value <= MAX_WALLET_LIMIT, "SafePump: Exceeds max wallet limit (1.5%)");
            }
        }

        // 4. Sell Limit check: max 2% total supply sold or transferred per 24h (except for excluded addresses)
        if (!isExcludedFromSellLimit[from]) {
            SellRecord storage record = sellRecords[from];
            if (block.timestamp - record.lastWindowStart >= 1 days) {
                record.lastWindowStart = uint128(block.timestamp);
                record.amountSold = uint128(value);
            } else {
                record.amountSold += uint128(value);
            }
            require(record.amountSold <= SELL_LIMIT_24H, "SafePump: Exceeds 24h sell limit (2%)");
        }

        super._update(from, to, value);
    }

    /**
     * @notice Calculate the amount of tokens currently locked for the creator
     * @return Locked token amount (in wei)
     */
    function getLockedCreatorAmount() public view returns (uint256) {
        if (!migrationComplete) {
            return CREATOR_ALLOCATION;
        }

        // Vesting starts 48 hours after migration
        uint256 vestingStart = migrationTime + 48 hours;
        if (block.timestamp < vestingStart) {
            return CREATOR_ALLOCATION;
        }

        uint256 timeElapsed = block.timestamp - vestingStart;
        uint256 daysElapsed = timeElapsed / 1 days;
        uint256 unlockedAmount = daysElapsed * VESTING_RELEASE_RATE_PER_DAY;

        if (unlockedAmount >= CREATOR_ALLOCATION) {
            return 0;
        }

        return CREATOR_ALLOCATION - unlockedAmount;
    }
}
