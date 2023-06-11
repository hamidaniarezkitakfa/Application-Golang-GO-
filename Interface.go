package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jlaffaye/ftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	// Boucle infinie pour afficher le menu à chaque fois que l'utilisateur a terminé une opération
	for {
		// Affichage des options disponibles
		fmt.Println("Sélectionnez une option:")
		fmt.Println("1. Créer un fichier")
		fmt.Println("2. Copier un fichier")
		fmt.Println("3. readFile")
		fmt.Println("4. Supprimer un fichier")
		fmt.Println("5. Afficher les permissions d'un fichier")
		fmt.Println("6. Se connecter à un serveur FTP")
		fmt.Println("7. Se connecter à un serveur HTTP")
		fmt.Println("8. Se connecter à un serveur SSH")
		fmt.Println("9. Quitter l'application")
		fmt.Println("10. Créer les tables de la base de données")
		fmt.Println("11. Ajouter un utilisateur")
		fmt.Println("12. Mettre à jour un utilisateur")
		fmt.Println("13. Supprimer un utilisateur")

		// Lecture du choix de l'utilisateur
		var choix int
		fmt.Scanln(&choix)
		// Exécution de l'action correspondante au choix de l'utilisateur
		switch choix {
		case 1:
			creerFichier()
		case 2:
			copierFichier()
		case 3:
			lireFichier()
		case 4:
			supprimerFichier()
		case 5:
			afficherPermissions()
		case 6:
			seConnecterFTP()
		case 7:
			seConnecterHTTP()
		case 8:
			seConnecterSSH()
		case 9:
			fmt.Println("Merci d'avoir utilisé notre application!")
			os.Exit(0)
		case 10:
			// créer les tables de la base de données
			err := createTables()
			if err != nil {
				fmt.Println("Erreur lors de la création des tables:", err)
			} else {
				fmt.Println("Les tables ont été créées avec succès!")
			}
		case 11:
			// ajouter un utilisateur
			var name, email string
			fmt.Print("Entrez le nom de l'utilisateur: ")
			fmt.Scan(&name)
			fmt.Print("Entrez l'email de l'utilisateur: ")
			fmt.Scan(&email)
			err := addUser(name, email)
			if err != nil {
				fmt.Println("Erreur lors de l'ajout de l'utilisateur:", err)
			}
		case 12:
			// mettre à jour un utilisateur
			var id int
			var name, email string
			fmt.Print("Entrez l'ID de l'utilisateur: ")
			fmt.Scan(&id)
			fmt.Print("Entrez le nouveau nom de l'utilisateur: ")
			fmt.Scan(&name)
			fmt.Print("Entrez le nouvel email de l'utilisateur: ")
			fmt.Scan(&email)
			err := updateUser(id, name, email)
			if err != nil {
				fmt.Println("Erreur lors de la mise à jour de l'utilisateur:", err)
			}
		case 13:
			// supprimer un utilisateur
			var id int
			fmt.Print("Entrez l'ID de l'utilisateur à supprimer: ")
			fmt.Scan(&id)
			err := deleteUser(id)
			if err != nil {
				fmt.Println("Erreur lors de la suppression de l'utilisateur:", err)
			}
		default:
			fmt.Println("Option invalide.")
		}
	}
}

func creerFichier() {
	// Saisire le nom du fichier
	fmt.Println("Entrez le nom du fichier que vous souhaitez créer : ")
	var nomFichier string
	fmt.Scanln(&nomFichier)

	// Crée un nouveau fichier avec le nom spécifié
	f, err := os.Create(nomFichier)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// Affiche un message de confirmation
	fmt.Println("Le fichier", nomFichier, "a été créé avec succès!")
}

func copierFichier() {
	// Lecture du nom du fichier source
	fmt.Println("Entrez le nom du fichier à copier : ")
	var source string
	fmt.Scanln(&source)

	// Lecture du nom du fichier de destination
	fmt.Println("Entrez le nom du fichier de destination : ")
	var destination string
	fmt.Scanln(&destination)

	// Lecture du contenu du fichier source
	input, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(destination, input, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Le fichier", source, "a été copié avec succès dans", destination)
}

func lireFichier() {
	// Lecture du nom du fichier à lire
	fmt.Print("Entrez le nom du fichier à lire : ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Erreur lors de la lecture de l'entrée utilisateur :", scanner.Err())
		return
	}
	nomFichier := scanner.Text()

	// Lecture du contenu du fichier
	content, err := ioutil.ReadFile(nomFichier)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}

	// Affichage du contenu du fichier
	fmt.Println("Contenu du fichier :")
	fmt.Println(string(content))
}

func supprimerFichier() {
	fmt.Println("Entrez le nom du fichier à supprimer : ")
	var nomFichier string
	fmt.Scanln(&nomFichier)
	//supprimer le fichier
	err := os.Remove(nomFichier)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Si la suppression a réussi, on affiche un message de confirmation
	fmt.Println("Le fichier", nomFichier, "a été supprimé avec succès!")
}

func afficherPermissions() {
	// Demande d'entrer le nom du fichier
	fmt.Println("Entrez le nom du fichier dont vous souhaitez afficher les permissions : ")
	var nomFichier string
	fmt.Scanln(&nomFichier)
	// Récupère les informations de statut du fichier
	f, err := os.Stat(nomFichier)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Affiche les permissions du fichier
	fmt.Println("Permissions de", nomFichier, ":", f.Mode().String())
}

func seConnecterFTP() {
	// Lecture des informations de connexion FTP
	fmt.Print("Entrez l'adresse du serveur FTP : ")
	var adresseServeur string
	fmt.Scanln(&adresseServeur)
	fmt.Print("Entrez le nom d'utilisateur : ")
	var nomUtilisateur string
	fmt.Scanln(&nomUtilisateur)
	fmt.Print("Entrez le mot de passe : ")
	var motDePasse string
	fmt.Scanln(&motDePasse)

	// Connexion au serveur FTP
	ftpClient, err := ftp.Connect(adresseServeur)
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur FTP :", err)
		return
	}
	defer ftpClient.Quit()

	err = ftpClient.Login(nomUtilisateur, motDePasse)
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur FTP :", err)
		return
	}

	// Récupération de la liste des fichiers sur le serveur FTP
	fichiers, err := ftpClient.NameList("/")
	if err != nil {
		log.Fatal(err)
	}

	// Affichage de la liste des fichiers
	fmt.Println("Liste des fichiers sur le serveur FTP :")
	for _, fichier := range fichiers {
		fmt.Println(fichier)
	}
}

func seConnecterHTTP() {
	// Entrer l'adresse du serveur HTTP
	fmt.Println("Entrez l'adresse du serveur HTTP : ")
	var adresse string
	fmt.Scanln(&adresse)
	// Effectue une requête GET vers l'adresse du serveur HTTP
	resp, err := http.Get(adresse)
	if err != nil {
		// Si une erreur se produit, l'affiche et retourne de la fonction
		fmt.Println(err)
		return
	}
	// Ferme le corps de la réponse de la requête après que la fonction a fini d'utiliser la réponse
	defer resp.Body.Close()

	fmt.Println("Connecté au serveur HTTP avec succès!")
}

func seConnecterSSH() {
	fmt.Println("Entrez l'adresse du serveur SSH : ")
	var adresse string
	fmt.Scanln(&adresse)

	fmt.Println("Entrez le nom d'utilisateur : ")
	var utilisateur string
	fmt.Scanln(&utilisateur)

	fmt.Println("Entrez le mot de passe : ")
	var motDePasse string
	fmt.Scanln(&motDePasse)
	// Configuration du client SSH
	config := &ssh.ClientConfig{
		User: utilisateur,
		Auth: []ssh.AuthMethod{
			ssh.Password(motDePasse),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connexion au serveur SSH
	client, err := ssh.Dial("tcp", adresse, config)
	if err != nil {
		fmt.Println("Erreur de connexion SSH :", err)
		return
	}
	// Ouverture d'une nouvelle session SSH
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Erreur de session SSH :", err)
		return
	}
	defer session.Close()

	fmt.Println("Connexion SSH réussie !")
}

// connexion à la base de données
func connectToDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/bddGo")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// fonction pour créer les tables si elles n'existent pas
func createTables() error {
	db, err := connectToDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// exécuter la requête de création de table pour chaque table nécessaire
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INT NOT NULL AUTO_INCREMENT,
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255) NOT NULL,
            PRIMARY KEY (id)
        )
    `)
	if err != nil {
		return err
	}

	// la fonction peut être étendue avec d'autres tables nécessaires à l'application

	return nil
}

// fonction pour ajouter un utilisateur
func addUser(name string, email string) error {
	db, err := connectToDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// exécuter la requête d'insertion d'utilisateur
	_, err = db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return err
	}

	fmt.Println("L'utilisateur", name, "a été ajouté avec succès!")
	return nil
}

// fonction pour mettre à jour un utilisateur
func updateUser(id int, name string, email string) error {
	db, err := connectToDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// exécuter la requête de mise à jour de l'utilisateur
	_, err = db.Exec("UPDATE users SET name=?, email=? WHERE id=?", name, email, id)
	if err != nil {
		return err
	}

	fmt.Println("L'utilisateur", name, "a été mis à jour avec succès!")
	return nil
}

// fonction pour supprimer un utilisateur
func deleteUser(id int) error {
	db, err := connectToDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// exécuter la requête de suppression de l'utilisateur
	_, err = db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return err
	}

	fmt.Println("L'utilisateur a été supprimé avec succès!")
	return nil
}
