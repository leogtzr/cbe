$(document).ready(function () {
    var persons = $('#persons');
    // REST call ... 
    if (persons.length) {
        console.log('Making http call ... ');

        $.ajax({
            url: '/persontypes',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    console.log("E: " + types[i]);
                }

            },
            error: function(data) {
                console.log('woops! :(' + data);
            }
        });

    }
});