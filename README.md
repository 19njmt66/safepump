# 🛡️ SafePump : Launchpad 100% Sécurisé (Base L2) - Fiche de Concept

Ce document centralise toutes les décisions de conception, l'analyse de marché, le modèle économique et l'architecture technique de notre future plateforme de lancement de tokens sécurisée, développée en Solidity sur le réseau Base.

---

## 1. 📊 Analyse de Marché & Opportunité
* **Le Problème :** Sur Pump.fun, **99,8% des tokens échouent** à migrer. La majorité absolue des lancements sont des arnaques directes ou des dumps rapides. Le fléau majeur est le *bundling au bloc 0* (le créateur achète 30% de la supply à bas coût via 20 wallets secrets et dump sur le public).
* **Le Positionnement :** Se positionner comme le **"Havre de Paix"** pour les acheteurs lassés des arnaques, et comme un **"Label de Confiance"** pour les créateurs honnêtes qui veulent prouver la légitimité de leur projet à leur communauté.

---

## 2. 🛡️ Modèle de Sécurité (La "Trust Architecture")
Le protocole applique des règles strictes inscrites en dur dans le smart contract de token ERC-20 :
1. **Anti-Honeypot :** Token standardisé ERC-20 immuable, sans propriétaire (*no owner*), sans liste noire (*no blacklist*), et sans taxe modifiable.
2. **Anti-LP Rugpull :** Les jetons de liquidité (LP) générés lors de la migration automatique sur Uniswap V2 sont **brûlés instantanément** dans la même transaction blockchain.
3. **Limite de Portefeuille Dynamique (Anti-Sybil & Anti-Bundle) :** 
   - *Phase de Lancement (10 premiers blocs) :* Un seul achat par bloc est autorisé à l'échelle globale du contrat + limite stricte de **2% de la supply max par portefeuille** (empêche le bundling à bas coût par le dev).
   - *Phase d'Incubation (0% à 20% de progression de la courbe) :* La limite de **2% max par portefeuille** reste active pour forcer une distribution saine de départ.
   - *Phase de Croissance (20% à 100% de progression) :* La limite est désactivée pour permettre aux whales et degens d'acheter de plus grosses positions une fois le lancement validé comme organique.
   - *Post-Migration (Uniswap V2) :* Toutes les limites de détention sont définitivement désactivées pour le trading libre.
4. **Avantage Créateur (Pack 5%) :** Le créateur du projet peut activer un bonus de **5% de la supply au lancement** en payant **0,02127 ETH** (environ 60 $). S'il ne l'active pas, il reçoit 0% de bonus gratuit, et le public achète ces 5% sur la courbe.
5. **Anti-Panic Dump (Sell Limit) :** Aucun portefeuille (hors DEX) ne peut vendre plus de **2% de la supply totale par tranche de 24h**, lissant les baisses et évitant les crashs à -99%.

---

## 💸 3. Modèle Économique & Incitations

### A. Pour le Créateur (Honnêteté Rentable)
* **Allocation de 5% Gratuite :** Ses jetons de créateur sont transparents, légitimes et verrouillés.
* **Bonus de Graduation :** Un bonus fixe en ETH (ex: 0.02 ETH) lui est versé automatiquement lors de la migration vers Uniswap.

### B. Pour la Plateforme (Tes Revenus)
* **100% des Frais de Trading :** Tu conserves l'intégralité (100%) des taxes de transactions de 1% collectées sur la courbe de liaison.
* **Frais de Migration :** Une taxe fixe (ex: 0.03 ETH) prélevée sur le pool lors de la graduation.
* **Espaces Publicitaires (Trending Slots) :** Vente d'emplacements sponsorisés sur la page d'accueil (ex: 0.05 ETH / jour).

---

## ⚙️ 4. Architecture Technique & Hébergement (Stratégie Hybride Offshore)

Pour offrir une vitesse d'exécution optimale tout en garantissant une résilience juridique absolue (Hors-MiCA, Hors-US), nous adoptons un modèle hybride de serveurs.

### A. Pile Technique Détaillée
* **Blockchain :** **Base (Layer 2 Ethereum)**. Frais de gaz infimes (< 0,01$) et support de Solidity pour le codage de nos limites de transfert personnalisées (anti-rug, vesting, cooldowns).
* **Nom de Domaine (DNS) :** Enregistrement anonyme payé en Monero (XMR) chez **Njalla** (extension `.is` pour l'Islande, réputée pour sa législation sur la liberté d'expression).
* **Frontend (Site Web) :** Développé en React/Next.js et hébergé de manière décentralisée sur le réseau **IPFS** via **Fleek**. Il est distribué en peer-to-peer mondial, le rendant techniquement incensurable (contrairement à un hébergement classique sur Vercel).
* **Backend API & Indexeur :** Écrit en **Go** (performances CPU élevées, idéal pour les WebSockets en direct) et hébergé sur un serveur **VPS offshore** situé en **Islande** (chez **FlokiNET**) ou en **Malaisie** (chez **Shinjiru**), payé anonymement en cryptomonnaie.
* **Base de données :** Base **PostgreSQL hébergée localement** sur le VPS offshore (nous excluons Supabase car ses serveurs sont situés aux USA et soumis aux lois américaines).
* **RPC Nodes :** Alchemy, QuickNode ou Ankr, configurés derrière un proxy/VPN sur le VPS pour masquer l'adresse IP du serveur.
* **Masquage IP & Sécurité :** Le serveur Go est protégé par un reverse-proxy offshore (comme **DDOS-Guard** ou Cloudflare) pour masquer son adresse IP réelle et le protéger des attaques.

### B. Étude de Cas : Infrastructure de Pump.fun
* **Leur choix :** Pump.fun utilise **Vercel** pour son site web, **AWS (Amazon)** pour ses bases de données et serveurs d'API, et **Cloudflare** pour ses DNS.
* **Leur stratégie :** Ils privilégient la puissance de serveurs corporate américains (très scalables) et se protègent uniquement par l'anonymat des fondateurs et des montages de sociétés écrans.
* **Notre amélioration :** Nous conservons la même efficacité technique (Base de données locale rapide + API Go), mais en décentralisant le site sur **IPFS** et en déplaçant les serveurs physiques dans des pays **hors MiCA et hors juridiction américaine (Islande/Malaisie)**. C'est le niveau maximal de sécurité et de conformité offshore.

### C. Onboarding Utilisateur & Portabilité (Coinbase SDK / Smart Wallet)
Pour attirer le grand public (Web2), la plateforme intègre le **Smart Wallet de Coinbase** (accès instantané par FaceID/Passkey sans phrase de récupération).
* **Compatibilité Blockchain (Permissionless) :** Le SDK de connexion au Smart Wallet est open-source et décentralisé. Coinbase fonctionne comme un protocole réseau : **ils ne peuvent pas bloquer notre dApp** ni censurer les connexions sur notre site offshore. L'intégration de la connexion wallet est 100% libre et sans autorisation préalable.
* **Le cas du Fiat Onramp (Achat direct par Carte/Apple Pay) :** L'intégration directe d'un widget d'achat par carte (ex: Coinbase Pay, Stripe Onramp) directement sur notre nom de domaine offshore présente un risque de refus ou de blocage par ces émetteurs réglementés.
* **Stratégie de contournement (Fallback) :**
  1. **Agrégateurs d'Onramp Flexibles :** Intégration de widgets tiers multi-fournisseurs (comme *Onramper* ou *Topper by Uphold*) qui sont beaucoup plus tolérants avec les plateformes de meme coins.
  2. **Achat natif dans le Wallet :** En cas de blocage du widget sur notre site, nous affichons un guide simple redirigeant l'utilisateur vers l'option d'achat *interne* de son application Coinbase Wallet (ce que Coinbase ne bloquera jamais car il s'agit de son propre client direct), puis le redirigeons sur notre site pour finaliser son swap en 1 clic.

---

## ⚖️ 5. Réglementation et Aspects Légaux

Lancer un launchpad de tokens implique des responsabilités juridiques majeures (lois sur les valeurs mobilières, régulations MiCA en Europe, règles anti-blanchiment AML/KYC).

### A. Les Risques Réglementaires Principaux
1. **Qualification en "Securities" (Titres Financiers) :** Aux USA (SEC) et dans d'autres pays, si un token est vendu avec une promesse de gain dépendant des efforts du créateur, il est considéré comme une valeur mobilière. Faciliter son échange sans licence expose la plateforme à des sanctions.
2. **Régulation MiCA (Europe) :** En Europe, les prestataires de services sur crypto-actifs (CASP) et les émetteurs de tokens doivent être enregistrés et publier un *whitepaper*. Les meme coins y sont de plus en plus surveillés.

### B. Stratégies de Protection Juridique (Mitigation)
Pour opérer la plateforme en toute sécurité sans subir de poursuites judiciaires, nous appliquerons 4 boucliers légaux :

* **1. Géo-blocage strict (Geoblocking IP) :**
  Bloquer l'accès à l'interface web pour les adresses IP issues des pays à haut risque réglementaire (États-Unis, Chine, etc.). C'est le standard de l'industrie DeFi pour prouver aux régulateurs qu'on ne cible pas leurs citoyens.
* **2. Modèle Non-Custodial (Purement Décentralisé) :**
  La plateforme ne détient **jamais** les fonds des utilisateurs. Tout est exécuté directement de portefeuille à portefeuille via des smart contracts autonomes sur la blockchain Base. C'est l'utilisateur qui signe et prend l'entière responsabilité technique de son swap.
* **3. Création d'une structure Offshore :**
  Enregistrer l'entité juridique qui exploite l'interface web de la plateforme dans une juridiction propice aux crypto-actifs (Panama, Îles Vierges Britanniques BVI, Seychelles ou Dubaï/Suisse) pour protéger les fondateurs à titre personnel.
* **4. Disclaimers & CGU Rigoureuses :**
  Faire accepter des termes d'utilisation stricts spécifiant en dur que :
  - Les jetons créés sont uniquement destinés au **divertissement (meme coins)**.
  - Ils n'ont **aucune valeur intrinsèque, aucune utilité financière, et aucune promesse de rendement**.
  - L'utilisateur assume 100% des risques de perte en capital.

---

## 🚀 Prochaines Étapes
1. Écriture du smart contract Solidity du token ERC-20 (`SafePumpToken.sol`) intégrant les règles de transfert (limites 2%/24h, 2% max wallet, block cooldown).
2. Développement du contrat Factory (`SafePumpFactory.sol`) qui déploie ces tokens, gère la courbe de liaison et effectue la migration automatique + burn de LP.
