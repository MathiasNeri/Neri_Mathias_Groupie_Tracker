
function searchAnime() {
    const searchBox = document.querySelector('.searchInput');
    const query = searchBox.value;
    if (query.trim()) {
        window.location.href = `/search?q=${encodeURIComponent(query)}`;
    } else {
        alert('Veuillez entrer un terme de recherche.');
    }
        
    console.log('toto');
}
