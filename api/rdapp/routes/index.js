var express = require('express');
var router = express.Router();
/* GET home page. */
router.get('/', function(req, res, next) {
  var cookie = req.cookies['session'];
  console.log("The cookie at index is " + cookie)
    if(cookie === undefined) {
      var title = req.query.title;
      if(title == '') {
        title = "RealDirect"
      }
      res.render('index', { title: title });
    } else {
      res.redirect('/dashboard?email=' + cookie);
    }
});

module.exports = router;
