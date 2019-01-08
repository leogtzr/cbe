$(document).ready(function () {

    $('#alert').hide();

    var persons = $('#persontypes');
    // REST call ... 
    if (persons.length) {
        $.ajax({
            url: '/persontypes',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    $('#persontypes').append($('<option name="' + types[i].ID + '">').append(types[i].Type));
                }
            },
            error: function(data) {
                console.log('woops! :(' + data);
            }
        });
    }

    $('#addperson').on('submit', function(e) {

        var currentForm = this;
        e.preventDefault();
        var name = $('#person_name').val();
        var personType = $('#persontypes').find(":selected").attr('name');

        console.log('Name: ' + name);
        console.log('Person type: ' + personType);

        $.ajax({
            url: '/addperson',
            type: 'POST',
            data: {name: name, type: personType},
            success: function(data) {
                $("#alert").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert").slideUp(500);
                });
            },
            error: function(data) {
                // TODO: Show alert ... 
                console.log(data);
            }
        });

    });

});