// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {SafePumpFactory} from "../src/SafePumpFactory.sol";

/**
 * @title DeploySafePump
 * @dev Deploy script for SafePumpFactory.
 * Usage:
 * forge script script/DeploySafePump.s.sol --rpc-url <rpc_url> --private-key <private_key> --broadcast
 */
contract DeploySafePump is Script {
    // Base Mainnet / Sepolia Uniswap V2 Router address
    // (Uniswap V2 Router is at the same address on both Base Mainnet and Base Sepolia)
    address public constant UNISWAP_V2_ROUTER = 0x4752ba5DBc23f44D87826276BF6Fd6b1C372aD24;

    function run() external {
        // Retrieve deployment configuration from environment variables or use defaults
        uint256 deployerPrivateKey = vm.envOr("DEPLOYER_PRIVATE_KEY", uint256(0));
        address feeRecipient = vm.envOr("FEE_RECIPIENT", address(0));
        address initialOwner = vm.envOr("INITIAL_OWNER", address(0));

        // If not specified in env, use default fallback for local/testnet testing
        if (deployerPrivateKey == 0) {
            console.log("WARNING: DEPLOYER_PRIVATE_KEY not set. Using default Anvil private key.");
            deployerPrivateKey = 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80;
        }

        address deployer = vm.addr(deployerPrivateKey);

        if (feeRecipient == address(0)) {
            feeRecipient = deployer;
        }

        if (initialOwner == address(0)) {
            initialOwner = deployer;
        }

        console.log("Deploying contracts with deployer:", deployer);
        console.log("Uniswap V2 Router:", UNISWAP_V2_ROUTER);
        console.log("Fee Recipient:", feeRecipient);
        console.log("Initial Owner:", initialOwner);

        vm.startBroadcast(deployerPrivateKey);

        SafePumpFactory factory = new SafePumpFactory(
            UNISWAP_V2_ROUTER,
            feeRecipient,
            initialOwner
        );

        vm.stopBroadcast();

        console.log("SafePumpFactory successfully deployed at:", address(factory));
    }
}
