function modifyFavorite(add, animeID) {
    const action = add ? 'addFavorite' : 'removeFavorite';
    const url = `/${action}/${animeID}`;

    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ animeID })
    })
        .then(response => {
            if (response.ok) {
                window.location.reload();
            } else {
                alert('Something went wrong');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}


document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll('.favorite-checkbox').forEach(function (checkbox) {
        updateFavoriteLabel(checkbox);
    });
});

function handleFavoriteChange(checkbox) {
    const animeID = checkbox.dataset.animeid;
    const addToFavorites = checkbox.checked;
    modifyFavorite(addToFavorites, animeID);
    updateFavoriteLabel(checkbox);
}

function updateFavoriteLabel(checkbox) {
    const isFavorite = checkbox.checked;
    const label = checkbox.nextElementSibling;
    label.querySelector('.option-add').style.display = isFavorite ? 'none' : 'block';
    label.querySelector('.option-remove').style.display = isFavorite ? 'block' : 'none';
}