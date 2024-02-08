$(document).ready(function() {
  $('#searchInput').on('keyup', function() {
    var query = $(this).val();
    if(query.length === 0){
      document.getElementById("searchResults").style.display = "none";
      document.getElementById("searchResults").innerHTML = "";
      return;
    }
    document.getElementById("searchResults").style.display = "block";
    $.ajax({
      url: '/?query='+query, // Update with the correct endpoint URL
      method: 'SEARCH', // Update with the desired HTTP method
      data: { query: query },
      success: function(response) {
        // Clear previous results
        $('#searchResults').empty();

        // Display new results
        var results = response;
        for (var i = 0; i < results.length; i++) {
          var result = results[i];
          if(result.Member !== ""){
            $('#searchResults').append('<a style="color:black;" class="sr" href="/viewartist?id='+result.ID+'"><div class="result" style="text-transform: uppercase;">' + result.Member + '<br><span style="font-size:smaller; color:grey;">'+result.Artist+' - Member</span></div></a>');
        } else if(result.Location !== ""){
          $('#searchResults').append('<a style="color:black;" class="sr" href="/viewartist?id='+result.ID+'"><div class="result" style="text-transform: uppercase;">' + result.Artist + '<br><span style="font-size:smaller; color:grey;">'+result.Location+'</span></div></a>');
        } else {
          $('#searchResults').append('<a style="color:black;" class="sr" href="/viewartist?id='+result.ID+'"><div class="result" style="text-transform: uppercase;">' + result.Artist +'</div></a>');
        }
        }
        
      },
      dataType: 'json'
    });
  });
});