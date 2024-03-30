
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

function applyFilters(event) {
    event.preventDefault();
    const form = document.getElementById('filterForm');
    const formData = new FormData(form);
    const types = formData.getAll('type'); // Récupère tous les types cochés
    const queryString = types.map(type => `type=${type}`).join('&');
    const searchQuery = document.getElementById('searchBox').value.trim();

    // Reconstruire l'URL avec les filtres et rediriger
    window.location.href = `/result_search?q=${encodeURIComponent(searchQuery)}&${queryString}`;
}
