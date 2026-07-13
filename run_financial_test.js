const { execSync } = require('child_process');

const rpcUrl = "http://127.0.0.1:8545";
const mnemonic = "test test test test test test test test test test test junk";
const factoryAddress = "0x9a676e781a523b5d0c0e43731313a708cb607508";

const addresses = {
  platform: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", // Account #0
  creator: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",  // Account #1
  traderA: "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",  // Account #2
  traderB: "0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1"   // Account #3
};

function getBalance(addr) {
  const balStr = execSync(`cast balance "${addr}" --ether --rpc-url "${rpcUrl}"`).toString().trim();
  return parseFloat(balStr);
}

function getTokenBalance(tokenAddr, userAddr) {
  try {
    const cmd = `cast call "${tokenAddr}" "balanceOf(address)(uint256)" "${userAddr}" --rpc-url "${rpcUrl}"`;
    const balRaw = execSync(cmd).toString().trim();
    const balHex = balRaw.split(' ')[0]; // Strip the bracket representation
    const balInt = BigInt(balHex);
    return Number(balInt) / 1e18;
  } catch (e) {
    console.error(`[DEBUG] getTokenBalance failed for token=${tokenAddr} user=${userAddr}:`, e.message);
    return 0;
  }
}

function setBalance(addr, ethAmount) {
  const wei = BigInt(ethAmount * 1e18);
  const hex = "0x" + wei.toString(16);
  execSync(`cast rpc anvil_setBalance "${addr}" "${hex}" --rpc-url "${rpcUrl}"`);
}

function getProgress(tokenAddr) {
  try {
    const soldRaw = execSync(`cast call "${factoryAddress}" "tokens(address)(address,address,uint256,uint256,bool)" "${tokenAddr}" --rpc-url "${rpcUrl}" | head -n 3 | tail -n 1`).toString().trim();
    const soldHex = soldRaw.split(' ')[0]; // Strip the bracket representation
    const sold = BigInt(soldHex);
    const percent = (Number(sold) / 800000000e18) * 100;
    return percent;
  } catch (e) {
    console.error(`[DEBUG] getProgress failed for token=${tokenAddr}:`, e.message);
    return 0;
  }
}

function getRaisedEth(tokenAddr) {
  try {
    const raisedRaw = execSync(`cast call "${factoryAddress}" "tokens(address)(address,address,uint256,uint256,bool)" "${tokenAddr}" --rpc-url "${rpcUrl}" | head -n 4 | tail -n 1`).toString().trim();
    const raisedHex = raisedRaw.split(' ')[0]; // Strip the bracket representation
    return Number(BigInt(raisedHex)) / 1e18;
  } catch (e) {
    console.error(`[DEBUG] getRaisedEth failed for token=${tokenAddr}:`, e.message);
    return 0;
  }
}

console.log("======================================================================");
console.log("🧪 SCÉNARIO DE SIMULATION FINANCIÈRE GLOBALE (3 WALLETS)");
console.log("======================================================================");

// 1. Réinitialisation des soldes à 10 ETH
console.log("\nReset des portefeuilles à 10.0 ETH...");
setBalance(addresses.platform, 10.0);
setBalance(addresses.creator, 10.0);
setBalance(addresses.traderA, 10.0);
setBalance(addresses.traderB, 10.0);

console.log(`Portefeuille Plateforme : ${getBalance(addresses.platform).toFixed(4)} ETH`);
console.log(`Portefeuille Créateur   : ${getBalance(addresses.creator).toFixed(4)} ETH`);
console.log(`Portefeuille Trader A   : ${getBalance(addresses.traderA).toFixed(4)} ETH`);
console.log(`Portefeuille Trader B   : ${getBalance(addresses.traderB).toFixed(4)} ETH`);

// 2. Création du Token par le Créateur
console.log("\n[ÉTAPE 1] Le Créateur lance le token 'SIMULATION' (SIM)...");
const createTxJson = execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index 1 "${factoryAddress}" "createToken(string,string)" "Simulation" "SIM" --json`).toString().trim();
const txData = JSON.parse(createTxJson);
const txHash = txData.transactionHash;

// Get receipt log to decode address
const receiptJson = execSync(`cast receipt --rpc-url "${rpcUrl}" "${txHash}" --json`).toString().trim();
const receipt = JSON.parse(receiptJson);
const log = receipt.logs.find(l => l.topics[0] === '0xd5d05a8421149c74fd223cfc823befb883babf9bf0b0e4d6bf9c8fdb70e59bb4');
const tokenAddress = '0x' + log.topics[1].slice(26);

console.log(`-> Token déployé avec succès à l'adresse : ${tokenAddress}`);

console.log("\nSoldes post-création :");
console.log(`- Créateur : ${getBalance(addresses.creator).toFixed(4)} ETH | MTK : ${getTokenBalance(tokenAddress, addresses.creator).toLocaleString()} SIM`);
console.log(`- Contrat Factory : MTK : ${getTokenBalance(tokenAddress, factoryAddress).toLocaleString()} SIM`);

// 3. Achat par le Trader A de 0.04 ETH
console.log("\n[ÉTAPE 2] Le Trader A achète pour 0.04 ETH...");
execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index 2 "${factoryAddress}" "buy(address)" "${tokenAddress}" --value 0.04ether`);

console.log("\nSoldes post-achat (0.04 ETH) :");
console.log(`- Plateforme : ${getBalance(addresses.platform).toFixed(4)} ETH (Attendu : ~10.00024 ETH)`);
console.log(`- Créateur   : ${getBalance(addresses.creator).toFixed(4)} ETH (Attendu : ~10.00016 ETH)`);
console.log(`- Trader A   : ${getBalance(addresses.traderA).toFixed(4)} ETH | MTK : ${getTokenBalance(tokenAddress, addresses.traderA).toLocaleString()} SIM`);
console.log(`- Factory Raised : ${getRaisedEth(tokenAddress).toFixed(4)} ETH (Attendu : 0.0396 ETH)`);

// 4. Vente partielle par le Trader A
console.log("\n[ÉTAPE 3] Le Trader A revend la moitié de ses jetons...");
const balanceA = getTokenBalance(tokenAddress, addresses.traderA);
const amountToSell = BigInt(Math.floor(balanceA / 2 * 1e18));

// Cooldown: Roll block
execSync(`cast rpc evm_mine --rpc-url "${rpcUrl}"`);

// Approve Factory
execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index 2 "${tokenAddress}" "approve(address,uint256)" "${factoryAddress}" ${amountToSell}`);
// Sell
execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index 2 "${factoryAddress}" "sell(address,uint256)" "${tokenAddress}" ${amountToSell}`);

console.log("\nSoldes post-vente :");
console.log(`- Plateforme : ${getBalance(addresses.platform).toFixed(5)} ETH`);
console.log(`- Créateur   : ${getBalance(addresses.creator).toFixed(5)} ETH`);
console.log(`- Trader A   : ${getBalance(addresses.traderA).toFixed(5)} ETH | MTK : ${getTokenBalance(tokenAddress, addresses.traderA).toLocaleString()} SIM`);

// 5. Test de la limite anti-whale (limite de 1.5% max wallet active pendant incubation)
console.log("\n[ÉTAPE 4] Le Trader B tente d'acheter pour 5 ETH (dépassement des 1.5% de supply)...");
try {
  execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index 3 "${factoryAddress}" "buy(address)" "${tokenAddress}" --value 5ether`, { stdio: 'ignore' });
  console.log("FAIL: La transaction d'achat Whale n'a pas été rejetée !");
} catch (e) {
  console.log("SUCCESS: L'achat Whale a été rejeté avec succès (limite max wallet active) !");
}

// 6. Complétion de la courbe de liaison (incubation puis gros achats)
console.log("\n[ÉTAPE 5] Phase d'incubation : 30 petits acheteurs de 0.03 ETH entrent pour franchir les 20% de progression...");
for (let index = 4; index <= 33; index++) {
  // Derive wallet address for the mnemonic index
  const addr = execSync(`cast wallet address --mnemonic "${mnemonic}" --mnemonic-index ${index}`).toString().trim();
  // Pre-fund the wallet with 1.0 ETH
  setBalance(addr, 1.0);

  execSync(`cast rpc evm_mine --rpc-url "${rpcUrl}"`);
  try {
    execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${index} "${factoryAddress}" "buy(address)" "${tokenAddress}" --value 0.03ether`);
  } catch (e) {
    console.error(`[ERROR] Achat échoué pour Wallet #${index} (Adresse: ${addr}). Détails :`);
    console.error(e.stderr ? e.stderr.toString() : e.message);
    throw e;
  }
}

console.log(`Progression de la courbe après incubation : ${getProgress(tokenAddress).toFixed(1)}%`);
console.log(`ETH levés après incubation : ${getRaisedEth(tokenAddress).toFixed(4)} ETH`);

// Maintenant que l'incubation est terminée (progression > 20%), la limite max wallet est désactivée.
// Nous pouvons effectuer de gros achats pour atteindre la limite de migration de 17 ETH.
console.log("\n[ÉTAPE 6] Limites désactivées. Le Trader B achète massivement pour pousser la courbe vers 100%...");
execSync(`cast rpc evm_mine --rpc-url "${rpcUrl}"`);
execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index 3 "${factoryAddress}" "buy(address)" "${tokenAddress}" --value 10ether`, { stdio: 'ignore' });

console.log("\n- Derniers achats pour finaliser les 17 ETH et déclencher la migration...");
for (let index = 34; index <= 38; index++) {
  try {
    const addr = execSync(`cast wallet address --mnemonic "${mnemonic}" --mnemonic-index ${index}`).toString().trim();
    setBalance(addr, 5.0);
    execSync(`cast rpc evm_mine --rpc-url "${rpcUrl}"`);
    execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index ${index} "${factoryAddress}" "buy(address)" "${tokenAddress}" --value 2.0ether`, { stdio: 'ignore' });
  } catch (e) {
    // Already migrated
  }
}

console.log(`\nProgression de la courbe : ${getProgress(tokenAddress).toFixed(1)}%`);
console.log(`ETH levés dans la Factory : ${getRaisedEth(tokenAddress).toFixed(4)} ETH`);

// Verification of migration status
const isMigrated = execSync(`cast call "${factoryAddress}" "tokens(address)(address,address,uint256,uint256,bool)" "${tokenAddress}" | tail -n 1`).toString().trim();
const migratedBool = isMigrated.includes("true") || isMigrated.includes("1") || isMigrated.endsWith("1");

if (migratedBool) {
  console.log("\n🎉 SUCCESS: Le token a franchi les 100% et a été migré avec succès sur Uniswap V2 !");
  console.log("\nSoldes finaux post-migration :");
  console.log(`- Plateforme : ${getBalance(addresses.platform).toFixed(4)} ETH (frais de trading + frais de migration 0.03 ETH inclus)`);
  console.log(`- Créateur   : ${getBalance(addresses.creator).toFixed(4)} ETH (frais de trading + bonus de migration 0.02 ETH inclus)`);
} else {
  console.log("\nLa migration n'a pas encore été déclenchée (solde insuffisant). Déclenchement forcé pour valider les calculs...");
  const addr9 = execSync(`cast wallet address --mnemonic "${mnemonic}" --mnemonic-index 9`).toString().trim();
  setBalance(addr9, 20.0);
  execSync(`cast rpc evm_mine --rpc-url "${rpcUrl}"`);
  execSync(`cast send --rpc-url "${rpcUrl}" --mnemonic "${mnemonic}" --mnemonic-index 9 "${factoryAddress}" "buy(address)" "${tokenAddress}" --value 10ether`, { stdio: 'ignore' });
  
  const finalMigrated = execSync(`cast call "${factoryAddress}" "tokens(address)(address,address,uint256,uint256,bool)" "${tokenAddress}" | tail -n 1`).toString().trim();
  const finalBool = finalMigrated.includes("1") || finalMigrated.endsWith("1");
  if (finalBool) {
    console.log("\n🎉 SUCCESS: Le token a franchi les 100% et a été migré avec succès sur Uniswap V2 !");
    console.log("\nSoldes finaux post-migration :");
    console.log(`- Plateforme : ${getBalance(addresses.platform).toFixed(4)} ETH`);
    console.log(`- Créateur   : ${getBalance(addresses.creator).toFixed(4)} ETH`);
  }
}

console.log("\n======================================================================");
console.log("FIN DU TEST D'INTÉGRATION FINANCIÈRE");
console.log("======================================================================");
