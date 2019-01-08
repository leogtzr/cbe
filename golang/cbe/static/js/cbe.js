$(document).ready(function () {
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
                    $('#persontypes').append($('<option>').append(types[i]));
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
        var personType = $('#persontypes').val();

        console.log('Hello!: ' + name + ", your type is: " + personType);

    });

});