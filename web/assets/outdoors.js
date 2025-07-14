document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementsByTagName('form')[0];
    const fields = Array.from(document.querySelectorAll('input, select, textarea')).reduce((acc, field) => {
        if (acc[field.name]) {
            if (Array.isArray(acc[field.name])) {
                acc[field.name].push(field);
            }
            else {
                acc[field.name] = [acc[field.name], field];
            }
        } else {
            acc[field.name] = field;
        }
        
        return acc;
    }, {});
    const errors = document.querySelectorAll('.error');

    // Handle showing optional fields
    Array.from(document.querySelectorAll('input[name="climbing_previous"], input[name="gear"]')).forEach(element => {
        const map = {
            'Bouldering': '#climbing_bouldering_field',
            'Rope': '#climbing_rope_field',
            'Harness': '#gear_harness_size_field',
            'Shoes': '#gear_shoe_size_field',
        };

        const onChange = () => {
            const selector = map[element.value];
            if (selector) {
                if (element.checked) {
                    document.querySelector(selector).classList.remove('hidden');
                } else {
                    document.querySelector(selector).classList.add('hidden');
                }
            }
        };
        
        element.addEventListener('change', onChange);
        onChange();
    });

    // Prevent multiple form submissions
    let submissionInProgress = false;

    // Handle form submission
    form.addEventListener('submit', async function(event) {
        event.preventDefault();
        
        let isValid = true;
        let firstError = null;

        // Exit early if a submission is already in progress
        if (submissionInProgress) {
            return;
        }

        // Reset errors
        errors.forEach(error => error.classList.add('hidden'));

        // Validate fields
        Object.values(fields).forEach(field => {
            if ((field.required || field?.dataset?.required === 'true') && !field.value.trim()) {
                document.getElementById(field.name + 'Error').classList.remove('hidden');
                if (!firstError) firstError = field;
                isValid = false;
            } else if (field.value && field.value.length > 0 && field.dataset.validation && !new RegExp(field.dataset.validation).test(field.value)) {
                document.getElementById(field.name + 'FormatError').classList.remove('hidden');
                if (!firstError) firstError = field;
                isValid = false;
            }
        });

        if (document.getElementById('gear_harness').checked && !document.getElementById('gear_harness_size').value.trim()) {
            document.getElementById('gear_harness_sizeError').classList.remove('hidden');
            if (!firstError) firstError = field;
            isValid = false;
        }

        if (document.getElementById('gear_shoes').checked && !document.getElementById('gear_shoe_size').value.trim()) {
            document.getElementById('gear_shoe_sizeError').classList.remove('hidden');
            if (!firstError) firstError = field;
            isValid = false;
        }

        if (!document.getElementById('disclaimer').checked) {
            document.getElementById('disclaimerError').classList.remove('hidden');
            if (!firstError) firstError = field;
            isValid = false;
        }

        if (isValid && !submissionInProgress) {
            grecaptcha.ready(function() {
                grecaptcha.execute(document.getElementById('recaptchaSiteKey').value, {action: 'submit'}).then(async function(token) {
                    submissionInProgress = true;
                    document.getElementById('submitButton').disabled = true;

                    const formSlug = document.getElementById('formSlug').value;
                    const response = await fetch(`/api/v1/submission/${formSlug}`, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            fname: fields['fname'].value,
                            lname: fields['lname'].value,
                            pronouns: fields['pronouns'].value,
                            email: fields['email'].value,
                            phone_number: fields['phone_number'].value,
                            ec_fname: fields['ec_fname'].value,
                            ec_lname: fields['ec_lname'].value,
                            ec_phone_number: fields['ec_phone_number'].value,
                            ec_secondary_phone_number: fields['ec_secondary_phone_number'].value,
                            ec_relation: fields['ec_relation'].value,
                            cp_riding: fields['cp_riding'].value,
                            cp_driving: fields['cp_driving'].value,
                            climbing_previous: fields['climbing_previous'].filter((e) => e.checked).map((e) => e.value).join(','),
                            climbing_bouldering: fields['climbing_previous'].filter((e) => e.checked && e.value === 'Bouldering').length === 1 ? fields['climbing_bouldering'].value : '',
                            climbing_rope: fields['climbing_previous'].filter((e) => e.checked && e.value === 'Rope').length === 1 ? fields['climbing_rope'].value : '',
                            climbing_goals: fields['climbing_goals'].value,
                            gear: fields['gear'].filter((e) => e.checked).map((e) => e.value).join(','),
                            gear_harness_size: fields['gear'].filter((e) => e.checked && e.value === 'Harness') ? fields['gear_harness_size'].value : '',
                            gear_shoe_size: fields['gear'].filter((e) => e.checked && e.value === 'Shoes') ? fields['gear_shoe_size'].value: '',
                            health_conditions: fields['health_conditions'].value,
                            health_accommodations: fields['health_accommodations'].value,
                            disclaimer: fields['disclaimer'].checked ? 'I have read and agree to the acknowledgment above.' : '',
                            recaptcha_token: token
                        })
                    });

                    if (response.status === 201) {
                        document.getElementById('successMessage').classList.remove('hidden');
                        form.classList.add('hidden');
                    } else {
                        const errorData = await response.json();
                        if (errorData.error) {
                            document.getElementById('errorMessage').innerText = 'An error occurred while submitting the form. Please try again. (' + response.status + ' - ' + errorData.error + ')';
                        } else {
                            document.getElementById('errorMessage').innerText = 'An error occurred while submitting the form. Please try again. (' + response.status + ')';
                        }
                        document.getElementById('errorMessage').classList.remove('hidden');
                    }
                
                    submissionInProgress = false;
                });
            });
        } else if (firstError) {
            // Focus on the first error field and scroll into the view.
            firstError.focus();
            firstError.scrollIntoView();
        }
    });
});