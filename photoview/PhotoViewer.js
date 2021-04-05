//Cognito_PhotoProcessingAppAuth_Role
//Cognito_PhotoProcessingAppUnauth_Role
//S3_Web_App_Access_Policy

// **DO THIS**:
//   Replace BUCKET_NAME with the bucket name.
//
var albumBucketName = 'myron--bucket';

// **DO THIS**:
//   Replace this block of code with the sample code located at:
//   Cognito -- Manage Identity Pools -- [identity_pool_name] -- Sample Code -- JavaScript
//
// Initialize the Amazon Cognito credentials provider
AWS.config.region = 'us-west-1'; // Region
AWS.config.credentials = new AWS.CognitoIdentityCredentials({
    IdentityPoolId: 'us-west-1:53f00295-3065-48c4-b86c-8b08949f280e',
});

// Create a new service object
var s3 = new AWS.S3({
  apiVersion: '2006-03-01',
  params: {Bucket: albumBucketName}
});

// A utility function to create HTML.
function getHtml(template) {
  return template.join('\n');
}

// List the photo albums that exist in the bucket.
function listAlbums() {
    s3.listObjects({Delimiter: '/'}, function(err, data) {
      if (err) {
        return alert('There was an error listing your albums: ' + err.message);
      } else {
          var albums = data.CommonPrefixes.map(function(commonPrefix) {
          var prefix = commonPrefix.Prefix;
          var albumName = decodeURIComponent(prefix.replace('/', ''));
          if(albumName.startsWith("album")){
            return getHtml([
              '<li>',
                '<button style="margin:5px;" onclick="viewAlbum(\'' + albumName + '\')">',
                  albumName,
                '</button>',
              '</li>'
            ]);
          }
          return ""
        });
        var message = albums.length ?
          getHtml([
            '<p>Click on an album name to view it.</p>',
          ]) :
          '<p>You do not have any albums. Please Create album.';
        var htmlTemplate = [
          '<h2>Albums</h2>',
          message,
          '<ul>',
            getHtml(albums),
          '</ul>',
        ]
        document.getElementById('viewer').innerHTML = getHtml(htmlTemplate);
      }
    });
  }


// Show the photos that exist in an album.
function viewAlbum(albumName) {
    var albumPhotosKey = encodeURIComponent(albumName) + '/';
    s3.listObjects({Prefix: albumPhotosKey}, function(err, data) {
      if (err) {
        return alert('There was an error viewing your album: ' + err.message);
      }
      // 'this' references the AWS.Request instance that represents the response
      var href = this.request.httpRequest.endpoint.href;
      var bucketUrl = href + albumBucketName + '/';
  
      var photos = data.Contents.map(function(photo) {
        var photoKey = photo.Key;
        var photoUrl = bucketUrl + encodeURIComponent(photoKey);


        var temp_name = photoKey.substr(photoKey.lastIndexOf('/')+1); 
        var only_name = temp_name.substring(0,temp_name.indexOf('.'));
        if(only_name.length != 0){
          return getHtml([
            '<span>',
              '<div>',
                '<br/>',
                '<img style="width:128px;height:128px;" src="' + photoUrl + '"/>',
              '</div>',
              '<div>',
                '<span>',
                  photoKey.replace(albumPhotosKey, ''),
                '</span>',
              '</div>',
              '<div>',
              '<button style="margin:5px;" onclick="getExif(\'' + only_name + '\'' + ',\'' + albumName + '\')">',
                'Show Exif Information',
              '</button>',
              '</div>',
            '</span>',
          ]);
        }
        return getHtml([
          '<span>',
            '<div>',
              '<br/>',
              '<img style="width:128px;height:128px;" src="' + photoUrl + '"/>',
            '</div>',
            '<div>',
              '<span>',
                photoKey.replace(albumPhotosKey, ''),
              '</span>',
            '</div>',
            '<div>',
            '</div>',
          '</span>',
        ]);
        
      });
      var message = photos.length ?
        '<p>The following photos are present.</p>' :
        '<p>There are no photos in this album.</p>';
      var htmlTemplate = [
        '<div>',
          '<button onclick="listAlbums()">',
            'Back To Albums',
          '</button>',
        '</div>',
        '<h2>',
          'Album: ' + albumName,
        '</h2>',
        message,
        '<div>',
          getHtml(photos),
        '</div>',
        '<h2>',
          'End of Album: ' + albumName,
        '</h2>',
        '<div>',
        '</div>',
      ]
      document.getElementById('viewer').innerHTML = getHtml(htmlTemplate);
      document.getElementsByTagName('img')[0].setAttribute('style', 'display:none;');
    });
  }

  function getExif(photoName, albumName) {
    var key = "exif/" + photoName + "_exif.txt"
    var params = {Bucket: 'myron--bucket', Key: key}
    s3.getObject(params, function(err, s3file) {
       var josnString = s3file.Body.toString('ascii');
       //exifJson = JSON.stringify(josnString)
       //console.log("exifJson:" + exifJson)
       var exif = JSON.parse(josnString);
       console.log(exif["Model"])
       console.log(exif["DateTime"])
       console.log(exif["GPSLatitude"])
       console.log(exif["GPSLongitude"])
       console.log(exif["GPSAltitude"])
       console.log(exif["ImageLength"])
       console.log(exif["ImageWidth"])
       
       html = getHtml([
        '<div>',
        '<button style="margin:5px;" onclick="viewAlbum(\'' + albumName + '\')">',
            'Back to ' + albumName,
        '</button>',
        '</div>',
        '<h2>EXIF Information</h2>',
         '<div>',
         '<span>',
         '<p>',
            'Model: ' + exif["Model"],
         '</p>',
         '<p>',
            'DateTime: ' + exif["DateTime"],
         '</p>',
         '<p>',
            'GPS Latitude: ' + exif["GPSLatitude"],
         '</p>',
         '<p>',
            'GPS Longitude: ' + exif["GPSLongitude"],
         '</p>',
         '<p>',
            'GPS Altitude: ' + exif["GPSAltitude"],
         '</p>',
         '<p>',
            'Image Length: ' + exif["ImageLength"],
         '</p>',
         '<p>',
            'Image Width: ' + exif["ImageWidth"],
         '</p>',

         '</span>',
         '</div>',
       ])
       document.getElementById('viewer').innerHTML = html

    })

    
  }