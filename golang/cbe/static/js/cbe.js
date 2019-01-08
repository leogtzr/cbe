$(document).ready(function () {
    var persons = $('#persons');
    // REST call ... 
    if (persons.length) {
        $.ajax({
            url: '/persontypes',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    $('#persons').append($('<option>').append(types[i]));
                }
            },
            error: function(data) {
                console.log('woops! :(' + data);
            }
        });
    }
});