package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"backend/config"
	"backend/database"

	"gorm.io/gorm"
)

// List of 10 standard Anvil accounts
var testWallets = []struct {
	Address  string
	Username string
	Bio      string
	Avatar   string
}{
	{"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "MemeKing 👑", "Le souverain suprême des dégénérés. Fondateur de multiples lancements.", "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&h=150&fit=crop"},
	{"0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "PumpMaster 💎", "Je ne vends jamais. Uniquement du diamant de qualité sur Base L2.", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop"},
	{"0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", "BasedChad 💪", "Sport, nutrition et trading de memecoins à fort levier.", "https://images.unsplash.com/photo-1620121692029-d088224ddc74?w=150&h=150&fit=crop"},
	{"0x90F79bf6EB2c4f870365E785982E1f101E93b906", "LunaLover 🌙", "Ancien astronaute reconverti dans la recherche de la prochaine gemme Base.", "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=150&h=150&fit=crop"},
	{"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65", "VitalicJr 🦄", "Je lis des contrats intelligents au petit-déjeuner. Code is law.", "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&h=150&fit=crop"},
	{"0x9965507D1a0565297baB5244DB855174956f4d90", "BaseGod ⚡", "Ici pour propager la bonne parole de Base L2. Constructeur et investisseur.", "https://images.unsplash.com/photo-1539571696357-5a69c17a67c6?w=150&h=150&fit=crop"},
	{"0x976EA74026E726554dB657fA54763abd0C3a0aa9", "DegenPro 🔥", "100% spéculatif. Je prends des risques démesurés pour le sport.", "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150&h=150&fit=crop"},
	{"0x14dC79964da2C08b2e2ec585030d1118165583f7", "PepeSpammer 🐸", "La grenouille est la seule unité de valeur qui compte.", "https://images.unsplash.com/photo-1522075469751-3a6694fb2f61?w=150&h=150&fit=crop"},
	{"0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f", "Dogefather 🐕", "Beaucoup de wow, très de gains. Le meilleur ami de l'homme sur la blockchain.", "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=150&h=150&fit=crop"},
	{"0xa0Ee7A142d267C1f36714E4a8F75612F20a79720", "AlphaCaller 📣", "Je partage mes analyses et mes pépites avant la foule. Suivez le guide.", "https://images.unsplash.com/photo-1580489944761-15a19d654956?w=150&h=150&fit=crop"},
}

// 50 unique memecoin definitions
var tokenTemplates = []struct {
	Name        string
	Symbol      string
	Description string
	Category    string
}{
	{"BasePaws", "PAWS", "L'amour des chiens sur Base L2. Rejoignez la meute et partagez des mèmes canins insolites !", "dog"},
	{"PepeFly", "PEFLY", "La grenouille verte qui voulait voler jusqu'à la lune grâce à ses petites ailes articulées.", "frog"},
	{"ChadDoge", "CDOGE", "Le Doge ultime qui passe tout son temps à la salle de sport à soulever de la fonte.", "dog"},
	{"MoonCat", "MCAT", "Un chat de l'espace curieux explorant les mystères de la galaxie et des trous noirs financiers.", "cat"},
	{"BaseGold", "BGOLD", "La réserve d'or virtuelle et hautement spéculative de la blockchain Base.", "gold"},
	{"FrogSuit", "FSUIT", "Une grenouille très sérieuse en costume trois pièces qui analyse les graphiques de la DeFi.", "frog"},
	{"DegenLife", "DEGEN", "Pour ceux qui ouvrent des positions avant de réfléchir. Bienvenue dans votre club privé.", "degen"},
	{"RocketBunny", "RBUN", "Un lapin survolté propulsé par une fusée d'énergie et de carottes virtuelles.", "rabbit"},
	{"CyberMeme", "CYBER", "Le futur des mèmes décentralisés dans un univers cyberpunk et rétro-éclairé au néon.", "cyber"},
	{"BaseShiba", "BSHIB", "La version ultra-rapide et sécurisée du célèbre Shiba Inu sur le réseau Base.", "dog"},
	{"BabyWojak", "BWOJ", "Le petit Wojak qui verse une larme dès que le cours du Bitcoin recule de 0.5% en une heure.", "wojak"},
	{"GigaChad", "GIGA", "Le jeton officiel de la confiance absolue et de la droiture d'esprit. Soyez un Chad.", "degen"},
	{"DogeLover", "DLVR", "Un hommage infini et affectueux à notre regretté Kabosu, le Doge originel.", "dog"},
	{"BasedFrog", "BFROG", "Une petite grenouille verte qui chill paisiblement au bord de l'eau sur le réseau Base.", "frog"},
	{"NeonCat", "NCAT", "Le chat légendaire qui parcourt le cosmos en laissant une traînée scintillante derrière lui.", "cat"},
	{"SpaceFrog", "SFROG", "Parce que l'espace infini a cruellement besoin de plus d'amphibiens courageux.", "frog"},
	{"BaseCoin", "BCOIN", "Le meme coin le plus simple, le plus brut et le plus efficace de toute la blockchain.", "degen"},
	{"TurboBase", "TURBO", "La vitesse supérieure et accélérée pour franchir la courbe de liaison le plus vite possible.", "degen"},
	{"GravityCoin", "GRV", "Le jeton physique qui prétend défier la gravité financière et les lois de l'attraction.", "moon"},
	{"SafeMoonBase", "SAFEM", "Le voyage vers la lune, mais cette fois-ci en toute sécurité et sans faille sur Base L2.", "moon"},
	{"PepeKing", "PEKING", "Le souverain incontesté et majestueux de toutes les grenouilles de l'Internet mondial.", "frog"},
	{"DogeKing", "DKING", "Le roi suprême qui gouverne avec sagesse le royaume des chiens rieurs.", "dog"},
	{"CatPower", "CATP", "Le pouvoir absolu et indomptable des félins domestiques appliqué à l'économie de la DeFi.", "cat"},
	{"BasePanda", "BPANDA", "Un panda zen qui mange du bambou en attendant sagement ses retours sur investissement.", "panda"},
	{"BasedHamster", "BHAM", "Le hamster trader équipé d'une roulette qui prend de meilleures décisions que les hedge funds.", "hamster"},
	{"MoonLover", "MLOVE", "Pour tous les rêveurs et romantiques qui passent leurs nuits les yeux rivés vers les étoiles.", "moon"},
	{"BaseRider", "RIDER", "Le motard de l'extrême qui pilote sa moto sans casque sur les courbes du marché.", "degen"},
	{"PixelDoge", "PDOGE", "Un Doge rétro dessiné en pixel-art 8 bits pour rappeler les consoles de notre enfance.", "dog"},
	{"BasedAlien", "ALIEN", "Une intelligence extraterrestre mystérieuse venue observer nos comportements financiers.", "cyber"},
	{"GoldenFrog", "GFROG", "La grenouille dorée légendaire qui apporte chance et prospérité à ses détenteurs.", "frog"},
	{"BaseNinja", "NINJA", "Le jeton furtif et discret qui franchit les étapes de liaison sans faire de bruit.", "degen"},
	{"MemeMachine", "MMACH", "La machine automatisée ultime conçue pour imprimer des mèmes drôles à la chaîne.", "cyber"},
	{"BaseViking", "VIK", "Des guerriers scandinaves robustes prêts à piller la liquidité et à conquérir la DeFi.", "degen"},
	{"SuperDoge", "SDOGE", "Le super-héros canin masqué dont le monde de la cryptomonnaie a désespérément besoin.", "dog"},
	{"SmartCat", "SCAT", "Un chat surdoué portant des lunettes qui passe ses journées à auditer le code des smart contracts.", "cat"},
	{"BaseWizard", "WIZ", "Le magicien barbu qui lance des sorts d'abondance sur les courbes de prix.", "cyber"},
	{"BabyPepe", "BPEPE", "La version bébé et innocente du grand Pepe, avec de grands yeux remplis d'espoir.", "frog"},
	{"BasedYeti", "YETI", "L'abominable homme des neiges quitte ses montagnes gelées pour s'installer sur Base.", "degen"},
	{"MoonShine", "SHINE", "La lueur argentée de la lune pour guider les investisseurs nocturnes à travers le brouillard.", "moon"},
	{"BasePirate", "PIRATE", "À l'abordage ! À la recherche du coffre au trésor caché dans la liquidité de Base L2.", "degen"},
	{"BasedKoala", "KOALA", "Le koala le plus détendu et paresseux de toute la jungle crypto. Zéro stress.", "panda"},
	{"NeonDoge", "NDOGE", "Un Doge stylisé illuminé par des tubes de néons fluorescents de couleur rose et bleue.", "dog"},
	{"BasedSloth", "SLOTH", "Le paresseux qui monte doucement vers les sommets, centime par centime.", "panda"},
	{"PepeGold", "PGOLD", "La monnaie royale en or frappée à l'effigie de notre cher Pepe l'amphibien.", "frog"},
	{"BaseFalcon", "FALCON", "Le faucon pèlerin qui surveille les mouvements du marché avec une précision chirurgicale.", "cyber"},
	{"MoonDegen", "MDEGEN", "Un astronaute audacieux prêt à risquer sa combinaison spatiale pour un aller simple.", "moon"},
	{"BasedTurtle", "TURTLE", "La tortue de la fable de La Fontaine : elle avance lentement mais finit par migrer.", "panda"},
	{"BaseGlitch", "GLITCH", "Une anomalie informatique visuelle très stylée dans la matrice des memecoins.", "cyber"},
	{"MoonWalker", "WALKER", "Un astronaute qui effectue son Moonwalk directement sur la courbe de graduation.", "moon"},
	{"BasedZeus", "ZEUS", "Le dieu de l'Olympe qui lance des éclairs de gains électriques sur ses fidèles.", "gold"},
}

// Deterministic mock image URLs based on category
var categoryImages = map[string][]string{
	"dog": {
		"https://images.unsplash.com/photo-1543466835-00a7907e9de1?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1537151608828-ea2b117b6f86?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1583511655857-d19b40a7a54e?w=300&h=300&fit=crop",
	},
	"cat": {
		"https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1533738363-b7f9aef128ce?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1573865526739-10659fec78a5?w=300&h=300&fit=crop",
	},
	"frog": {
		"https://images.unsplash.com/photo-1579357776188-c5ab044d759f?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1545063877-c1e82ab903ec?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1590005354167-6da97870c913?w=300&h=300&fit=crop",
	},
	"panda": {
		"https://images.unsplash.com/photo-1564349683136-77e08dba1ef7?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1526336024174-e58f5cdd8e13?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1615087240969-eeff2fa558f2?w=300&h=300&fit=crop",
	},
	"gold": {
		"https://images.unsplash.com/photo-1610375228911-c4ab6410213b?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=300&h=300&fit=crop",
	},
	"cyber": {
		"https://images.unsplash.com/photo-1509198397868-475647b2a1e5?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1485827404703-89b55fcc595e?w=300&h=300&fit=crop",
	},
	"moon": {
		"https://images.unsplash.com/photo-1506703719100-a0f3a48c0f86?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1541185933-ef5d8ed016c2?w=300&h=300&fit=crop",
	},
	"degen": {
		"https://images.unsplash.com/photo-1621761191319-c6fb62004040?w=300&h=300&fit=crop",
		"https://images.unsplash.com/photo-1605792657660-596af9009e82?w=300&h=300&fit=crop",
	},
}

func main() {
	log.Println("=== DÉBUT DU SEMAGE DE LA BASE DE DONNÉES (50 JETONS) ===")

	// 1. Charge la configuration
	cfg := config.LoadConfig()

	// 2. Initialise la connexion GORM
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Échec de connexion à la base : %v", err)
	}

	// 3. Purge des tables existantes
	log.Println("Purge des tables en cours...")
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&database.Trade{}).Error; err != nil {
		log.Fatalf("Erreur purge Trades : %v", err)
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&database.Token{}).Error; err != nil {
		log.Fatalf("Erreur purge Tokens : %v", err)
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&database.User{}).Error; err != nil {
		log.Fatalf("Erreur purge Users : %v", err)
	}
	log.Println("Purge terminée.")

	// 4. Création des 10 profils d'utilisateurs test
	log.Println("Création des 10 profils d'utilisateurs test...")
	for _, wallet := range testWallets {
		user := database.User{
			Address:   strings.ToLower(wallet.Address),
			Username:  wallet.Username,
			Bio:       wallet.Bio,
			AvatarUrl: wallet.Avatar,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Erreur lors de la création de l'utilisateur %s : %v", wallet.Username, err)
		}
	}
	log.Println("Profils utilisateurs créés.")

	// Source de nombres aléatoires pour simuler la cohérence
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 5. Génération des 50 jetons
	log.Println("Création des 50 jetons et de leurs historiques de transactions...")
	
	// Limite de courbe de liaison complète : 800 000 000 de jetons vendus (correspond à 5 ETH levés)
	maxSupplyToSell := 800000000.0

	for i, t := range tokenTemplates {
		// Génère une adresse de contrat fictive unique pour le jeton
		tokenAddr := fmt.Sprintf("0x%040x", i+1)

		// Sélectionne un créateur aléatoire parmi nos 10 wallets
		creatorWallet := testWallets[r.Intn(len(testWallets))]
		creatorAddr := strings.ToLower(creatorWallet.Address)

		// Liste d'IDs Picsum garantis actifs et sans doublons
		picsumIDs := []int{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
			41, 42, 43, 44, 45, 48, 49, 50, 51, 52,
		}
		imgUrl := fmt.Sprintf("https://picsum.photos/id/%d/300/300", picsumIDs[i])

		// Choisit le degré de progression
		var tokensSold float64
		var ethRaised float64
		migrated := false
		pairAddress := ""
		hasCreatorBonus := r.Float64() < 0.4 // 40% de chance d'avoir activé le pack créateur

		// Répartition des états :
		// - 10 jetons migrés (diplômés / liquidité verrouillée à 100% sur Uniswap)
		// - 10 jetons très avancés (progress > 75%)
		// - 15 jetons moyennement avancés (progress entre 20% et 75%)
		// - 15 jetons récents (progress < 20%)
		if i < 10 {
			// MIGRÉS
			migrated = true
			tokensSold = maxSupplyToSell
			ethRaised = 5.0
			pairAddress = fmt.Sprintf("0x%040x", i+500) // paire Uniswap fictive
		} else if i < 20 {
			// PROGRESSION ÉLEVÉE
			progressPct := 0.75 + r.Float64()*0.23 // 75% à 98%
			tokensSold = maxSupplyToSell * progressPct
			ethRaised = 5.0 * progressPct
		} else if i < 35 {
			// PROGRESSION MOYENNE
			progressPct := 0.20 + r.Float64()*0.55 // 20% à 75%
			tokensSold = maxSupplyToSell * progressPct
			ethRaised = 5.0 * progressPct
		} else {
			// RÉCENTS
			progressPct := 0.0 + r.Float64()*0.20 // 0% à 20%
			tokensSold = maxSupplyToSell * progressPct
			ethRaised = 5.0 * progressPct
		}

		// Optionnel : Bonus créateur 5% (50M jetons) déjà soustraits ou inclus
		// Dans le contrat, le pack créateur coûte >= 0.02127 ETH
		if hasCreatorBonus && ethRaised < 0.02127 {
			ethRaised = 0.02127
			tokensSold = 50000000.0
		}

		// Dates cohérentes dans le passé
		createdAt := time.Now().Add(-time.Hour * time.Duration((50-i)*8)) // échelonné sur les 16 derniers jours

		// Formate les grands nombres en chaîne décimale (18 décimales)
		tokensSoldStr := fmt.Sprintf("%.0f000000000000000000", tokensSold)
		ethRaisedStr := fmt.Sprintf("%.0f", ethRaised*1e18)

		// Liens sociaux cohérents
		website := fmt.Sprintf("https://www.%s.io", strings.ToLower(t.Symbol))
		twitter := fmt.Sprintf("https://x.com/%s_coin", strings.ToLower(t.Symbol))
		telegram := fmt.Sprintf("https://t.me/%s_chat", strings.ToLower(t.Symbol))

		// Enregistre le jeton
		token := database.Token{
			Address:     tokenAddr,
			Creator:     creatorAddr,
			Name:        t.Name,
			Symbol:      t.Symbol,
			TokensSold:  tokensSoldStr,
			EthRaised:   ethRaisedStr,
			Migrated:    migrated,
			PairAddress: pairAddress,
			Description: t.Description,
			ImageUrl:    imgUrl,
			Website:     website,
			Twitter:     twitter,
			Telegram:    telegram,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		}

		if err := db.Create(&token).Error; err != nil {
			log.Fatalf("Erreur lors de la création du jeton %s : %v", t.Name, err)
		}

		// 6. Génération des Transactions (Trades) pour chaque jeton de manière COHÉRENTE
		// Si progress > 0, on génère plusieurs transactions d'achat qui somment exactement au total
		if tokensSold > 0 {
			numTrades := 3 + r.Intn(8) // 3 à 10 transactions d'achat par jeton
			
			// Répartition du total
			tokensRemaining := tokensSold
			ethRemaining := ethRaised * 1e18

			for k := 0; k < numTrades; k++ {
				var tokenAmt float64
				var ethAmt float64

				if k == numTrades-1 {
					// Dernier trade : prend le reste exact pour garantir la cohérence mathématique
					tokenAmt = tokensRemaining
					ethAmt = ethRemaining
				} else {
					// Trade intermédiaire : fraction aléatoire (entre 10% et 40% du restant)
					fraction := 0.10 + r.Float64()*0.30
					tokenAmt = tokensRemaining * fraction
					ethAmt = ethRemaining * fraction

					// Arrondit
					tokenAmt = float64(int64(tokenAmt))
					ethAmt = float64(int64(ethAmt))

					tokensRemaining -= tokenAmt
					ethRemaining -= ethAmt
				}

				if tokenAmt <= 0 {
					continue
				}

				// Acheteur aléatoire
				buyerWallet := testWallets[r.Intn(len(testWallets))]
				buyerAddr := strings.ToLower(buyerWallet.Address)

				// Tx hash fictif
				txHash := fmt.Sprintf("0x%064x", r.Uint64())

				// Date de transaction échelonnée entre le lancement et maintenant
				tradeTime := createdAt.Add(time.Minute * time.Duration(k*30))

				trade := database.Trade{
					TokenAddress:  tokenAddr,
					TxHash:        txHash,
					BlockNumber:   1000000 + uint64(i*1000+k),
					Timestamp:     tradeTime,
					IsBuy:         true,
					BuyerOrSeller: buyerAddr,
					TokenAmount:   fmt.Sprintf("%.0f000000000000000000", tokenAmt),
					EthAmount:     fmt.Sprintf("%.0f", ethAmt),
					Fee:           fmt.Sprintf("%.0f", ethAmt*0.01), // 1% de frais
				}

				if err := db.Create(&trade).Error; err != nil {
					log.Fatalf("Erreur lors de la création du trade pour %s : %v", t.Symbol, err)
				}
			}
		}
	}

	log.Println("50 jetons et leurs historiques de transactions créés avec succès !")
	log.Println("=== FIN DU SEMAGE DE LA BASE DE DONNÉES ===")
}
