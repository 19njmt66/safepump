/**
 * SafePump Zero-Dependency Activity & Volume Bot
 * 
 * Simulates realistic user activity on SafePump by:
 * - Launching new tokens (with or without creator packs)
 * - Buying random amounts of active tokens
 * - Selling portions of holdings
 * 
 * Uses Foundry's `cast` utility under the hood. No npm installation needed!
 * 
 * Usage:
 * 1. Run: node activity_bot.js
 */

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

// Load configuration from backend/.env if available
let rpcUrl = "http://127.0.0.1:8545";
let factoryAddress = "";
let mnemonic = "test test test test test test test test test test test junk";
let apiUrl = "http://localhost:8080/api/v1";

try {
  const envPath = path.join(__dirname, "backend", ".env");
  if (fs.existsSync(envPath)) {
    const content = fs.readFileSync(envPath, "utf-8");
    
    const rpcMatch = content.match(/RPC_URL=(.+)/);
    if (rpcMatch) rpcUrl = rpcMatch[1].trim();

    const factoryMatch = content.match(/FACTORY_ADDRESS=(0x[a-fA-F0-9]{40})/);
    if (factoryMatch) factoryAddress = factoryMatch[1].trim();

    const serverMatch = content.match(/SERVER_ADDR=(.+)/);
    if (serverMatch) {
      const port = serverMatch[1].replace(":", "").trim();
      apiUrl = `http://localhost:${port}/api/v1`;
    }
  }
} catch (e) {
  console.log("Could not load backend/.env, using default local config.");
}

if (!factoryAddress) {
  // Try to load factory address from config API
  try {
    const configRaw = execSync(`curl -s ${apiUrl}/config`).toString().trim();
    const config = JSON.parse(configRaw);
    factoryAddress = config.factory_address;
  } catch (e) {
    console.error("❌ Error: FACTORY_ADDRESS not found in backend config or API.");
    console.error("Please ensure the backend and anvil are running before starting the bot.");
    process.exit(1);
  }
}

// Bot logic parameters
const MIN_DELAY_MS = 6000;  // 6s
const MAX_DELAY_MS = 20000; // 20s
const CREATE_PROBABILITY = 0.05; // 5% chance to deploy a new token
const BUY_PROBABILITY = 0.60;   // 60% chance to buy
const SELL_PROBABILITY = 0.35;  // 35% chance to sell

const MIN_BUY_ETH = 0.005;
const MAX_BUY_ETH = 0.04;

// Token name generators
const adjectives = [
  "Based", "Pepe", "Degen", "Turbo", "Giga", "Mochi", "Toshi", "Roost", "Brett", "Chomp", 
  "Andy", "Mog", "Popcat", "Slerf", "Spitfire", "Milady", "Nyan", "Doge", "Shiba", "Floki", 
  "Bonk", "Myro", "Wojak", "Pudgy", "Smol", "Fren", "Cyber", "Rocket", "Space", "Neon",
  "Glow", "Lucky", "Happy", "Chad", "Quantum", "Hyper", "Alpha", "Based", "Coinbase", "Base"
];

const nouns = [
  "Brett", "Pepe", "Toshi", "Mochi", "Roost", "Chomp", "Frog", "Cat", "Dog", "Wolf", 
  "Panda", "Base", "Coin", "Token", "Meme", "Shib", "Inu", "Pudgy", "Penguin", "Mogger", 
  "Milady", "Wojak", "Chad", "Nyan", "Keyboard", "GigaChad", "Moon", "Toad", "Banana", "Hamster"
];

function generateRandomName() {
  const adj = adjectives[Math.floor(Math.random() * adjectives.length)];
  const noun = nouns[Math.floor(Math.random() * nouns.length)];
  const name = Math.random() > 0.5 ? `${adj} ${noun}` : `${noun} ${adj}`;
  const symbol = `${adj.slice(0, 3).toUpperCase()}${noun.slice(0, 3).toUpperCase()}`;
  return { name, symbol };
}

function getWalletAddress(index) {
  return execSync(`cast wallet address --mnemonic "${mnemonic}" --mnemonic-index ${index}`).toString().trim();
}

function getEthBalance(walletIndex) {
  const addr = getWalletAddress(walletIndex);
  const raw = execSync(`cast balance "${addr}" --rpc-url "${rpcUrl}"`).toString().trim();
  return BigInt(raw);
}

function getTokenBalance(tokenAddr, walletAddress) {
  try {
    const cmd = `cast call "${tokenAddr}" "balanceOf(address)(uint256)" "${walletAddress}" --rpc-url "${rpcUrl}"`;
    const balRaw = execSync(cmd).toString().trim();
    const balHex = balRaw.split(' ')[0];
    return BigInt(balHex);
  } catch (e) {
    return 0n;
  }
}

function checkAllowance(tokenAddr, walletAddress) {
  try {
    const cmd = `cast call "${tokenAddr}" "allowance(address,address)(uint256)" "${walletAddress}" "${factoryAddress}" --rpc-url "${rpcUrl}"`;
    const allowRaw = execSync(cmd).toString().trim();
    const allowHex = allowRaw.split(' ')[0];
    return BigInt(allowHex);
  } catch (e) {
    return 0n;
  }
}

function approveFactory(tokenAddr, walletIndex) {
  try {
    const cmd = `cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${walletIndex} "${tokenAddr}" "approve(address,uint256)" "${factoryAddress}" "115792089237316195423570985008687907853269984665640564039457584007913129639935" --json`;
    execSync(cmd, { stdio: "ignore" });
    return true;
  } catch (e) {
    return false;
  }
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function start() {
  console.log(`==================================================`);
  console.log(`🤖 SAFEPUMP SIMULATOR & ACTIVITY BOT ACTIVE`);
  console.log(`🔗 RPC Node : ${rpcUrl}`);
  console.log(`🏫 Factory  : ${factoryAddress}`);
  console.log(`📡 API URL  : ${apiUrl}`);
  console.log(`==================================================\n`);

  while (true) {
    // Pick a random wallet index between 0 and 9
    const walletIndex = Math.floor(Math.random() * 10);
    const walletAddress = getWalletAddress(walletIndex);

    console.log(`--------------------------------------------------`);
    console.log(`Wallet Selected: ${walletAddress} (Index: ${walletIndex})`);

    try {
      const ethBal = getEthBalance(walletIndex);
      console.log(`ETH Balance: ${(Number(ethBal) / 1e18).toFixed(4)} ETH`);

      if (ethBal < 50000000000000000n) { // < 0.05 ETH
        if (rpcUrl.includes("127.0.0.1") || rpcUrl.includes("localhost")) {
          console.log(`Balance low, refunding via Anvil RPC...`);
          // Refunding from anvil account
          execSync(`cast rpc anvil_setBalance "${walletAddress}" "0x3635c9adc5dea00000" --rpc-url "${rpcUrl}"`); // Set to 1000 ETH
          console.log(`Refunded to 1000 ETH.`);
        } else {
          console.log(`⚠️ Solde insuffisant ! Veuillez alimenter le portefeuille ${walletAddress} en ETH.`);
          continue; // Skip this wallet for this round
        }
      }

      // Fetch active tokens from API
      let activeTokens = [];
      try {
        const tokensRaw = execSync(`curl -s ${apiUrl}/tokens`).toString().trim();
        const tokens = JSON.parse(tokensRaw);
        activeTokens = tokens.filter(t => !t.migrated);
      } catch (e) {
        console.log("Could not fetch tokens list from API, trying to query contract directly...");
        // Fallback: Query allTokens from factory
        try {
          const allTokensRaw = execSync(`cast call "${factoryAddress}" "getAllTokens()(address[])" --rpc-url "${rpcUrl}"`).toString().trim();
          const addrs = allTokensRaw.replace(/[\[\]]/g, "").split(",").map(a => a.trim()).filter(a => a.startsWith("0x"));
          activeTokens = addrs.map(addr => ({ address: addr }));
        } catch (err) {
          console.error("Failed to query token list.");
        }
      }

      console.log(`Active tokens count: ${activeTokens.length}`);

      // Choose action
      let action = "buy";
      const rand = Math.random();

      if (activeTokens.length === 0 || rand < CREATE_PROBABILITY) {
        action = "create";
      } else if (rand < CREATE_PROBABILITY + BUY_PROBABILITY) {
        action = "buy";
      } else {
        action = "sell";
      }

      // Execute Action
      if (action === "create") {
        const { name, symbol } = generateRandomName();
        console.log(`🤖 Action: CREATE TOKEN [${name} (${symbol})]`);

        const buyPack = Math.random() < 0.3; // 30% chance to buy creator pack
        const valueParam = buyPack ? " --value 0.02127ether" : "";

        console.log(`Deploying token (Pack active: ${buyPack})...`);
        const cmd = `cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${walletIndex} "${factoryAddress}" "createToken(string,string)" "${name}" "${symbol}"${valueParam} --json`;
        
        const res = execSync(cmd).toString().trim();
        const receipt = JSON.parse(res);
        console.log(`Token created! Tx: ${receipt.transactionHash}`);

      } else if (action === "buy") {
        const targetTokenObj = activeTokens[Math.floor(Math.random() * activeTokens.length)];
        const targetToken = targetTokenObj.address;
        const ethToSpend = (MIN_BUY_ETH + Math.random() * (MAX_BUY_ETH - MIN_BUY_ETH)).toFixed(4);

        // Check if the token is still in incubation (< 20% progress)
        const tokensSold = targetTokenObj.tokens_sold ? BigInt(targetTokenObj.tokens_sold) : 0n;
        const isIncubation = tokensSold < BigInt("160000000000000000000000000");

        if (isIncubation) {
          // Estimate tokens to be received from this buy
          const Tv = BigInt("960000000000000000000000000"); // 960M
          const Ev = BigInt("1000000000000000000"); // 1 ETH
          const raised = targetTokenObj.eth_raised ? BigInt(targetTokenObj.eth_raised) : 0n;
          const currentTokenReserves = Tv - tokensSold;
          const currentEthReserves = Ev + raised;
          
          const ethIn = BigInt(Math.floor(parseFloat(ethToSpend) * 1e18));
          const netEth = (ethIn * 99n) / 100n;
          const numerator = currentTokenReserves * netEth;
          const denominator = currentEthReserves + netEth;
          const tokensToReceive = numerator / denominator;

          const tokenBalance = getTokenBalance(targetToken, walletAddress);
          const limit = BigInt("20000000000000000000000000"); // 2% of supply

          if (tokenBalance + tokensToReceive > limit) {
            console.log(`Buy of ${ethToSpend} ETH would exceed 2% wallet limit. Skipping buy.`);
            continue;
          }
        }
        
        console.log(`🤖 Action: BUY Token ${targetToken} for ${ethToSpend} ETH`);

        const cmd = `cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${walletIndex} "${factoryAddress}" "buy(address)" "${targetToken}" --value ${ethToSpend}ether --json`;
        const res = execSync(cmd).toString().trim();
        const receipt = JSON.parse(res);
        console.log(`Bought successfully! Tx: ${receipt.transactionHash}`);

      } else if (action === "sell") {
        const targetToken = activeTokens[Math.floor(Math.random() * activeTokens.length)].address;
        const walletAddress = getWalletAddress(walletIndex);
        const tokenBalance = getTokenBalance(targetToken, walletAddress);

        console.log(`Token Balance: ${(Number(tokenBalance) / 1e18).toFixed(0)} tokens`);

        if (tokenBalance === 0n) {
          console.log(`No tokens to sell. Skipping...`);
          continue;
        }

        const sellPct = 0.20 + Math.random() * 0.80; // Sell 20% to 100%
        const sellAmount = (tokenBalance * BigInt(Math.floor(sellPct * 100))) / 100n;

        console.log(`🤖 Action: SELL ${(sellPct * 100).toFixed(0)}% of Token holdings (${(Number(sellAmount) / 1e18).toFixed(0)} tokens)`);

        // Check allowance
        const allowance = checkAllowance(targetToken, walletAddress);
        if (allowance < sellAmount) {
          console.log(`Approving factory...`);
          approveFactory(targetToken, walletIndex);
          console.log(`Allowance approved.`);
        }

        const cmd = `cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${walletIndex} "${factoryAddress}" "sell(address,uint256)" "${targetToken}" "${sellAmount}" --json`;
        const res = execSync(cmd).toString().trim();
        const receipt = JSON.parse(res);
        console.log(`Sold successfully! Tx: ${receipt.transactionHash}`);
      }

    } catch (e) {
      console.log(`❌ Action failed: ${e.message}`);
    }

    const delay = Math.floor(MIN_DELAY_MS + Math.random() * (MAX_DELAY_MS - MIN_DELAY_MS));
    console.log(`Waiting ${(delay / 1000).toFixed(1)}s before next action...\n`);
    await sleep(delay);
  }
}

start().catch(e => {
  console.error("Fatal bot error:", e);
});
