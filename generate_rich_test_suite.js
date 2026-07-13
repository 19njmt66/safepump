const { execSync } = require('child_process');

const rpcUrl = "http://127.0.0.1:8545";
const mnemonic = "test test test test test test test test test test test junk";

// Get active factory address from backend config
let factoryAddress = "";
try {
  const configRaw = execSync('curl -s http://localhost:8080/api/v1/config').toString().trim();
  const config = JSON.parse(configRaw);
  factoryAddress = config.factory_address;
  console.log(`Using active factory address: ${factoryAddress}`);
} catch (e) {
  console.error("Failed to fetch factory address from backend config API:", e.message);
  process.exit(1);
}

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

function setBalance(addr, ethAmount) {
  const wei = BigInt(ethAmount * 1e18);
  const hex = "0x" + wei.toString(16);
  execSync(`cast rpc anvil_setBalance "${addr}" "${hex}" --rpc-url "${rpcUrl}"`);
}

function getWalletAddress(index) {
  return execSync(`cast wallet address --mnemonic "${mnemonic}" --mnemonic-index ${index}`).toString().trim();
}

function getTokenBalance(tokenAddr, userAddr) {
  try {
    const cmd = `cast call "${tokenAddr}" "balanceOf(address)(uint256)" "${userAddr}" --rpc-url "${rpcUrl}"`;
    const balRaw = execSync(cmd).toString().trim();
    const balHex = balRaw.split(' ')[0];
    const balInt = BigInt(balHex);
    return balInt;
  } catch (e) {
    return 0n;
  }
}

// Generate unique token names
const generatedNames = new Set();
function generateNameSymbol(i) {
  let name = "";
  let symbol = "";
  while (true) {
    const adj = adjectives[Math.floor(Math.random() * adjectives.length)];
    const noun = nouns[Math.floor(Math.random() * nouns.length)];
    name = Math.random() > 0.5 ? `${adj} ${noun}` : `${noun} ${adj}`;
    symbol = `${adj.slice(0,3).toUpperCase()}${noun.slice(0,3).toUpperCase()}`;
    if (generatedNames.has(name)) {
      name = `${name} ${i}`;
      symbol = `${symbol}${i}`;
    }
    if (!generatedNames.has(name)) {
      generatedNames.add(name);
      return { name, symbol };
    }
  }
}

function deployToken(name, symbol, walletIndex, initialValueEth = 0) {
  try {
    let cmd = `cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${walletIndex} "${factoryAddress}" "createToken(string,string)" "${name}" "${symbol}" --json`;
    if (initialValueEth > 0) {
      cmd += ` --value ${initialValueEth}ether`;
    }
    const createTxJson = execSync(cmd).toString().trim();
    const txData = JSON.parse(createTxJson);
    const txHash = txData.transactionHash;

    const receiptJson = execSync(`cast receipt --rpc-url "${rpcUrl}" "${txHash}" --json`).toString().trim();
    const receipt = JSON.parse(receiptJson);
    const log = receipt.logs.find(l => l.topics[0] === '0xd5d05a8421149c74fd223cfc823befb883babf9bf0b0e4d6bf9c8fdb70e59bb4');
    const tokenAddress = '0x' + log.topics[1].slice(26);
    return tokenAddress;
  } catch (e) {
    console.error(`Failed to deploy ${name}:`, e.message);
    return null;
  }
}

function buyToken(tokenAddress, walletIndex, amountEth) {
  try {
    execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${walletIndex} "${factoryAddress}" "buy(address)" "${tokenAddress}" --value ${amountEth}ether --json`, { stdio: 'ignore' });
    return true;
  } catch (e) {
    console.error(`Failed buy for token ${tokenAddress} from wallet ${walletIndex} of ${amountEth} ETH:`, e.message);
    return false;
  }
}

function sellToken(tokenAddress, walletIndex, tokenAmount) {
  try {
    // 1. Approve
    execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${walletIndex} "${tokenAddress}" "approve(address,uint256)" "${factoryAddress}" ${tokenAmount} --json`, { stdio: 'ignore' });
    // 2. Sell
    execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${walletIndex} "${factoryAddress}" "sell(address,uint256)" "${tokenAddress}" ${tokenAmount} --json`, { stdio: 'ignore' });
    return true;
  } catch (e) {
    console.error(`Failed sell for token ${tokenAddress} from wallet ${walletIndex}:`, e.message);
    return false;
  }
}

// Bypasses 2% max wallet limit by executing 12 small buys of 0.018 ETH from separate wallets
function pushToIncubationComplete(tokenAddress) {
  for (let w = 2; w <= 13; w++) {
    buyToken(tokenAddress, w, 0.018);
  }
}

async function main() {
  console.log("Pre-funding all wallets with 1000 ETH...");
  const walletAddresses = [];
  for (let i = 1; i <= 25; i++) {
    const addr = getWalletAddress(i);
    walletAddresses[i] = addr;
    setBalance(addr, 1000.0);
  }
  console.log("All wallets funded.");

  const countGroup1 = 15; // 0% progress
  const countGroup2 = 25; // Progress & active trades
  const countGroup3 = 10; // Migrated (100% progress)

  let tokenCount = 1;

  console.log(`\n--- GROUP 1: Deploying ${countGroup1} tokens with 0% progress ---`);
  for (let i = 0; i < countGroup1; i++) {
    const { name, symbol } = generateNameSymbol(tokenCount);
    const creatorIndex = i < 5 ? 1 : 2 + (i % 8);
    // Keep initial buy below 2% limit (0.015 ETH gets ~15M tokens)
    const initialBuy = (i % 3 === 0) ? 0.015 : 0;
    const address = deployToken(name, symbol, creatorIndex, initialBuy);
    if (address) {
      console.log(`[${tokenCount}/50] Deployed ${name} (${symbol}) at ${address} (0% progress)`);
      tokenCount++;
    }
  }

  console.log(`\n--- GROUP 2: Deploying ${countGroup2} tokens with active trading scenarios ---`);
  
  // Scenario 2A: PumpMaster Winning (10 tokens)
  console.log("-> Scenario 2A: 10 Tokens where PumpMaster is winning (unrealized profit)");
  for (let i = 0; i < 10; i++) {
    const { name, symbol } = generateNameSymbol(tokenCount);
    const address = deployToken(name, symbol, 2);
    if (address) {
      // PumpMaster buys early (under 2% max wallet limit)
      buyToken(address, 1, 0.015);
      // Others buy later, pushing price up
      buyToken(address, 3, 0.015);
      buyToken(address, 4, 0.015);
      buyToken(address, 5, 0.015);
      console.log(`[${tokenCount}/50] Deployed ${name} (${symbol}) at ${address} - PumpMaster Winning`);
      tokenCount++;
    }
  }

  // Scenario 2B: PumpMaster Losing (10 tokens)
  console.log("-> Scenario 2B: 10 Tokens where PumpMaster is losing (unrealized loss)");
  for (let i = 0; i < 10; i++) {
    const { name, symbol } = generateNameSymbol(tokenCount);
    const address = deployToken(name, symbol, 2);
    if (address) {
      // Wallet 3 buys first
      buyToken(address, 3, 0.015);
      // PumpMaster buys high
      buyToken(address, 1, 0.015);
      // Wallet 3 sells all its tokens to drop the price
      const bal3 = getTokenBalance(address, walletAddresses[3]);
      if (bal3 > 0n) {
        sellToken(address, 3, bal3);
      }
      console.log(`[${tokenCount}/50] Deployed ${name} (${symbol}) at ${address} - PumpMaster Losing`);
      tokenCount++;
    }
  }

  // Scenario 2C: PumpMaster fully exited (5 tokens)
  console.log("-> Scenario 2C: 5 Tokens where PumpMaster has fully exited (realized PnL)");
  for (let i = 0; i < 5; i++) {
    const { name, symbol } = generateNameSymbol(tokenCount);
    const address = deployToken(name, symbol, 2);
    if (address) {
      const isProfit = (i % 2 === 0);
      if (isProfit) {
        // Profit: Buy early, others buy, sell
        buyToken(address, 1, 0.015);
        buyToken(address, 3, 0.015);
        buyToken(address, 4, 0.015);
        const pmBal = getTokenBalance(address, walletAddresses[1]);
        if (pmBal > 0n) {
          sellToken(address, 1, pmBal);
        }
      } else {
        // Loss: Buy high, others dump, sell
        buyToken(address, 3, 0.015);
        buyToken(address, 1, 0.015);
        const bal3 = getTokenBalance(address, walletAddresses[3]);
        if (bal3 > 0n) {
          sellToken(address, 3, bal3);
        }
        const pmBal = getTokenBalance(address, walletAddresses[1]);
        if (pmBal > 0n) {
          sellToken(address, 1, pmBal);
        }
      }
      console.log(`[${tokenCount}/50] Deployed ${name} (${symbol}) at ${address} - PumpMaster Exited (${isProfit ? 'Profit' : 'Loss'})`);
      tokenCount++;
    }
  }

  console.log(`\n--- GROUP 3: Deploying ${countGroup3} fully migrated tokens (100% progress) ---`);
  for (let i = 0; i < countGroup3; i++) {
    const { name, symbol } = generateNameSymbol(tokenCount);
    const address = deployToken(name, symbol, 2);
    if (address) {
      // 1. Disable max wallet limit
      pushToIncubationComplete(address);
      // 2. Buy heavily from wallet 9 to trigger migration (5 ETH is enough to complete target)
      buyToken(address, 9, 5.0);
      console.log(`[${tokenCount}/50] Deployed and migrated ${name} (${symbol}) at ${address}`);
      tokenCount++;
    }
  }

  console.log("\n======================================================================");
  console.log("🎉 SUCCESS: Rich 50-token test suite created successfully!");
  console.log("======================================================================");
}

main().catch(console.error);
