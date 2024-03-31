const animeCards = document.querySelectorAll('.card');
animeCards.forEach(card => {
    card.addEventListener('click', function() {
        const animeID = card.getAttribute('data-anime-id');
        if (animeID) {
            const detailURL = `/anime_detail/${animeID}`;
            window.location.href = detailURL;
        }
    });
});
