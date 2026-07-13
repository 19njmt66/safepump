#!/usr/bin/env bash

# Find active local network IP address
LOCAL_IP=$(ip route get 1.1.1.1 2>/dev/null | grep -oP 'src \K\S+' || hostname -I | awk '{print $1}')

echo "======================================================================"
echo "🚀 DÉMARRAGE DES SERVICES SAFEPUMP ET PARTAGE RÉSEAU"
echo "🔗 IP de votre machine sur le réseau local : $LOCAL_IP"
echo "======================================================================"

# 1. Start Postgres database container
echo "📦 [1/4] Démarrage de PostgreSQL via Docker Compose..."
docker compose up -d

# 2. Start Anvil blockchain node exposed to network
echo "🔗 [2/4] Démarrage d'Anvil (blockchain de test) sur toutes les interfaces (0.0.0.0)..."
anvil --host 0.0.0.0 > anvil.log 2>&1 &
ANVIL_PID=$!

# Wait for Anvil to be ready
echo -n "   Attente du démarrage d'Anvil..."
until (echo > /dev/tcp/127.0.0.1/8545) >/dev/null 2>&1; do
  echo -n "."
  sleep 0.5
done
echo " OK!"

# 3. Deploy contracts to local Anvil
echo "📜 [3/4] Déploiement des smart contracts sur Anvil..."
DEPLOY_OUT=$(forge script script/DeploySafePumpLocal.s.sol:DeploySafePumpLocal --rpc-url http://127.0.0.1:8545 --broadcast)

# Extract deployed addresses from stdout
FACTORY_ADDR=$(echo "$DEPLOY_OUT" | grep "LOCAL_FACTORY_ADDRESS:" | awk '{print $2}')
ROUTER_ADDR=$(echo "$DEPLOY_OUT" | grep "MOCK_UNISWAP_V2_ROUTER:" | awk '{print $2}')

if [ -z "$FACTORY_ADDR" ]; then
  echo "❌ Erreur : Le déploiement des contrats a échoué. Arrêt d'Anvil."
  kill $ANVIL_PID
  exit 1
fi

echo "   -> Factory déployé à : $FACTORY_ADDR"
echo "   -> Router déployé à  : $ROUTER_ADDR"

# Create/Overwrite backend config file (.env)
echo "📝 Mise à jour de la configuration du backend (backend/.env)..."
cat <<EOF > backend/.env
DATABASE_URL=host=localhost user=postgres password=postgres dbname=safepump port=5432 sslmode=disable
RPC_URL=http://localhost:8545
FACTORY_ADDRESS=$FACTORY_ADDR
START_BLOCK=1
SERVER_ADDR=:8080
EOF

# 4. Start Go Backend in background
echo "⚙️ [4/4] Démarrage du backend Go (Indexer + API)..."
cd backend
go run main.go > backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# 5. Start React Frontend (Vite) exposed on the network
echo "💻 Démarrage du serveur de développement Frontend Vite sur 0.0.0.0:5173..."
cd frontend
npm run dev -- --host --port 5173 > frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

echo "======================================================================"
echo "🎉 TOUS LES SERVICES SONT EN COURS D'EXÉCUTION !"
echo "======================================================================"
echo "👉 Accès local (votre machine) :"
echo "   - Application Web : http://localhost:5173"
echo "   - Blockchain RPC  : http://localhost:8545"
echo "   - API Backend     : http://localhost:8080/api/v1"
echo ""
echo "👉 Accès réseau local (depuis un autre PC / smartphone sur le même Wi-Fi) :"
echo "   - Application Web : http://$LOCAL_IP:5173"
echo "   - Blockchain RPC  : http://$LOCAL_IP:8545"
echo "   - API Backend     : http://$LOCAL_IP:8080/api/v1"
echo "======================================================================"
echo "⚠️  Appuyez sur [ENTRÉE] pour arrêter tous les services..."
read -r

echo "🛑 Arrêt des services en cours..."
kill $FRONTEND_PID 2>/dev/null
kill $BACKEND_PID 2>/dev/null
kill $ANVIL_PID 2>/dev/null
docker compose down
echo "👋 Tous les services sont arrêtés !"
