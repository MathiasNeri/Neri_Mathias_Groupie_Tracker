# Groupie-Tracker
 
1-Lancement de la solution :

2-Ouvrez un terminal sur votre machine.
3-Naviguez jusqu'au répertoire où se trouve le code source de l'application.
4-Assurez-vous que Go est installé sur votre système.
5-Exécutez la commande go run . pour démarrer le serveur.
6- Ouvrez votre navigateur et accédez à http://localhost:8080/index pour afficher la page d'accueil.


Listing des routes :

1-GET /index
Fonctionnalité : Affiche la page d'accueil avec une sélection aléatoire d'animes et les animes à venir.

2-GET /result_search?q={query}&page={page}&type={type}
Arguments :

query : terme de recherche.
page : numéro de la page pour la pagination.
type : filtre par type d'anime.
Fonctionnalité : Effectue une recherche d'animes selon le terme donné, avec pagination et filtres.
3-GET /genres/
Fonctionnalité : Affiche les différentes catégories de genres d'animes.

4-GET /animes_by_genre/{genreID}?page={page}
Arguments :

genreID : identifiant unique du genre.
page : numéro de la page pour la pagination.
Fonctionnalité : Affiche les animes appartenant à un genre spécifique avec la possibilité de paginer.
5-POST /addFavorite/{animeID}
Arguments :

animeID : identifiant unique de l'anime à ajouter aux favoris.
Fonctionnalité : Ajoute un anime à la liste des favoris.
6-POST /removeFavorite/{animeID}
Arguments :

animeID : identifiant unique de l'anime à retirer des favoris.
Fonctionnalité : Retire un anime de la liste des favoris.
7-GET /favorites
Fonctionnalité : Affiche la liste des animes favoris de l'utilisateur.

8-GET /a_propos
Fonctionnalité : Affiche la page "À propos" avec des informations sur l'application et le développeur.