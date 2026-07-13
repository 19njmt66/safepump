const { execSync } = require('child_process');

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

const factoryAddress = "0x9a676e781a523b5d0c0e43731313a708cb607508";
const rpcUrl = "http://127.0.0.1:8545";
const mnemonic = "test test test test test test test test test test test junk";

console.log("Starting deployment of 100 credible meme coins from non-creator wallets...");

const generatedNames = new Set();
const deployed = [];

for (let i = 1; i <= 100; i++) {
  // Generate a unique name
  let name = "";
  let symbol = "";
  let attempts = 0;
  
  while (attempts < 100) {
    const adj = adjectives[Math.floor(Math.random() * adjectives.length)];
    const noun = nouns[Math.floor(Math.random() * nouns.length)];
    
    // Sometimes swap name order for diversity
    name = Math.random() > 0.5 ? `${adj} ${noun}` : `${noun} ${adj}`;
    symbol = `${adj.slice(0,3).toUpperCase()}${noun.slice(0,3).toUpperCase()}`;
    
    // Add a number suffix only if there is a collision
    if (generatedNames.has(name)) {
      name = `${name} ${i}`;
      symbol = `${symbol}${i}`;
    }
    
    if (!generatedNames.has(name)) {
      generatedNames.add(name);
      break;
    }
    attempts++;
  }

  // Cycle through mnemonic indices 1 to 9 (connected wallet is index 0)
  const mnemonicIndex = 1 + (i % 9); 

  console.log(`[${i}/100] Deploying ${name} (${symbol}) using Wallet #${mnemonicIndex}...`);
  try {
    const cmd = `cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${mnemonicIndex} ${factoryAddress} "createToken(string,string)" "${name}" "${symbol}"`;
    execSync(cmd, { stdio: 'ignore' });
    deployed.push({ name, symbol });
  } catch (e) {
    console.error(`Failed to deploy ${name}:`, e.message);
  }
}

console.log(`Successfully deployed ${deployed.length} tokens!`);
