document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('mnccForm');
    const errors = document.querySelectorAll('.error');

    const shoeSizeContainer = document.getElementById('shoeSizeContainer');
    const shoeSize = document.getElementById('shoeSize');
    const shoeSizeError = document.getElementById('shoeSizeError');


    // Show/hide shoe size based on selection
    const shoes = document.getElementById('shoes');
    shoeSizeContainer.classList.toggle('hidden', !shoes.checked);
    shoes.addEventListener('change', function() {
        shoeSizeContainer.classList.toggle('hidden', !this.checked);
    });

    form.addEventListener('submit', function(e) {
        e.preventDefault();
        
        let isValid = true;
        errors.forEach(error => error.classList.add('hidden'));

        // Validate required fields
        ['name', 'phoneNumber', 'email'].forEach(field => {
            const input = document.getElementById(field);
            const error = document.getElementById(field + 'Error');
            if (!input.value.trim()) {
                error.classList.remove('hidden');
                isValid = false;
            }
        });

        // Validate shoe size if climbing shoes are needed
        if (shoes.checked && !shoeSize.value.trim()) {
            shoeSizeError.classList.remove('hidden');
            isValid = false;
        }

        if (isValid) {
            // Here you would typically send the form data to your server
            console.log('Form submitted:', new FormData(form));
            alert('Form submitted successfully!');
            // Uncomment the next line to reset the form after submission
            // form.reset();
        }
    });
});