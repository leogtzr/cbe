$(document).ready(function () {

    $('#alert').hide();
    $('#alert_error').hide();

    var persontypes = $('#persontypes');
    if (persontypes.length) {
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

    var persons = $('#persons');
    if (persons.length) {
        $.ajax({
            url: '/persons',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    $('#persons').append($('<option name="' + types[i].ID + '">').append(types[i].Name));
                }
            },
            error: function(data) {
                console.log('woops! :(' + data);
            }
        });
    }

    var familyInteractions = $('#family_interactions');
    if (familyInteractions.length) {
        $.ajax({
            url: '/personspertype/1',
            type: 'GET',
            data: {},
            success: function(data) {
                var types = JSON.parse(data);
                for (i = 0; i < types.length; i++) {
                    $('#family_interactions').append($('<li name="1" class="list-group-item">').append(types[i]));
                    console.log("Persona: " + types[i]);
                }
            },
            error: function(data) {
                console.log(data)
                console.log('woops! :(' + data + ", not able to get types");
            }
        });
    }

    $('#addperson').on('submit', function(e) {

        var currentForm = this;
        e.preventDefault();
        var name = $('#person_name').val();
        var personType = $('#persontypes').find(":selected").attr('name');

        $.ajax({
            url: '/addperson',
            type: 'POST',
            data: {name: name, type: personType},
            success: function(data) {
                console.log("Good");
                $('#person_name').val('');
                $("#alert").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert").slideUp(500);
                });
            },
            error: function(data) {
                console.log("Error!");
                $("#alert_error").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert_error").slideUp(500);
                });
            }
        });

    });

    $('#addinteraction').on('submit', function(e) {

        var currentForm = this;
        e.preventDefault();
        var text = $('#interactiontext').val();
        var personId = $('#persons').find(":selected").attr('name');

        $.ajax({
            url: '/addinteraction',
            type: 'POST',
            data: {personId: personId, comment: text},
            success: function(data) {
                text.val('');
                $("#alert").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert").slideUp(500);
                });
            },
            error: function(data) {
                $("#alert_error").fadeTo(2000, 500).slideUp(500, function() {
                    $("#alert_error").slideUp(500);
                });
            }
        });

    });

});