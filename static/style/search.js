
document.addEventListener('DOMContentLoaded', (event) => {
    document.getElementById('searchForm').addEventListener('submit', function(e) {
        e.preventDefault(); // Empêche le formulaire de soumettre de manière traditionnelle
        searchAnime();
    });
});

function searchAnime() {
    const searchBox = document.querySelector('.searchInput');
    const query = searchBox.value;
    if (query.trim()) {
        window.location.href = `/result_search?q=${encodeURIComponent(query)}`;
    } else {
        alert('Veuillez entrer un terme de recherche.');
    }
}

