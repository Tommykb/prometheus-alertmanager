'use strict';

angular.module('am.directives', []);

angular.module('am.directives').directive('route',
  function(RecursionHelper) {
    return {
      restrict: 'E',
      scope: {
        route: '='
      },
      templateUrl: 'app/partials/route.html',
      compile: function(element) {
        // Use the compile function from the RecursionHelper,
        // And return the linking function(s) which it returns
        return RecursionHelper.compile(element);
      }
    };
  }
);

angular.module('am.directives').directive('alert',
  function() {
    return {
      restrict: 'E',
      scope: {
        alert: '=',
        group: '='
      },
      templateUrl: 'app/partials/alert.html'
    };
  }
);

angular.module('am.directives').directive('silence',
  function() {
    return {
      restrict: 'E',
      scope: {
        sil: '='
      },
      templateUrl: 'app/partials/silence.html'
    };
  }
);

angular.module('am.directives').directive('silenceForm',
  function() {
    return {
      restrict: 'E',
      scope: {
        silence: '=',
        formTitle: '@',
        formSubtitle: '@',
        buttonText: '@'
      },
      templateUrl: 'app/partials/silence-form.html'
    };
  }
);

angular.module('am.services', ['ngResource']);

angular.module('am.services').factory('Silence',
  function($resource) {
    return $resource('', {
      id: '@id'
    }, {
      'query': {
        method: 'GET',
        url: 'api/v1/silences'
      },
      'create': {
        method: 'POST',
        url: 'api/v1/silences'
      },
      'get': {
        method: 'GET',
        url: 'api/v1/silence/:id'
      },
      'delete': {
        method: 'DELETE',
        url: 'api/v1/silence/:id'
      }
    });
  }
);

angular.module('am.services').factory('Alert',
  function($resource) {
    return $resource('', {}, {
      'query': {
        method: 'GET',
        url: 'api/v1/alerts'
      }
    });
  }
);

angular.module('am.services').factory('AlertGroups',
  function($resource) {
    return $resource('', {}, {
      'query': {
        method: 'GET',
        url: 'api/v1/alerts/groups'
      }
    });
  }
);

angular.module('am.services').factory('Alert',
  function($resource) {
    return $resource('', {}, {
      'query': {
        method: 'GET',
        url: 'api/v1/alerts'
      }
    });
  }
);
angular.module('am.controllers', []);

angular.module('am.controllers').controller('NavCtrl',
  function($scope, $location) {
    $scope.items = [{
      name: 'Silences',
      url: 'silences'
    }, {
      name: 'Alerts',
      url: 'alerts'
    }, {
      name: 'Status',
      url: 'status'
    }];

    $scope.selected = function(item) {
      return item.url == $location.path()
    }
  }
);

angular.module('am.controllers').controller('AlertCtrl',
  function($scope) {
    $scope.showDetails = false;

    $scope.toggleDetails = function() {
      $scope.showDetails = !$scope.showDetails
    }

    $scope.showSilenceForm = false;

    $scope.toggleSilenceForm = function() {
      $scope.showSilenceForm = !$scope.showSilenceForm
    }

    $scope.silence = {
      matchers: []
    }
    angular.forEach($scope.alert.labels, function(value, key) {
      this.push({
        name: key,
        value: value,
        isRegex: false
      });
    }, $scope.silence.matchers);

    $scope.$on('silence-created', function(evt) {
      $scope.toggleSilenceForm();
    });
  }
);

angular.module('am.controllers').controller('AlertsCtrl',
  function($scope, $location, AlertGroups) {
    $scope.groups = null;
    $scope.allReceivers = [];

    $scope.$watch('receivers', function(recvs) {
      if (recvs === undefined || angular.equals(recvs, $scope.allReceivers)) {
        return;
      }
      if (recvs) {
        $location.search('receiver', recvs);
      } else {
        $location.search('receiver', null);
      }
    });

    $scope.notEmpty = function(group) {
      var ret = false

      angular.forEach(group.blocks, function(blk) {
        if (this.indexOf(blk.routeOpts.receiver) >= 0) {
          var unsilencedAlerts = blk.alerts.filter(function (a) { return !a.silenced; });
          if (!$scope.hideSilenced && blk.alerts.length > 0 || $scope.hideSilenced && unsilencedAlerts.length > 0) {
            ret = true
          }
        }
      }, $scope.receivers);

      return ret;
    };

    $scope.refresh = function() {
      AlertGroups.query({},
        function(data) {
          $scope.groups = data.data;

          $scope.allReceivers = [];
          angular.forEach($scope.groups, function(group) {
            angular.forEach(group.blocks, function(blk) {
              if (this.indexOf(blk.routeOpts.receiver) < 0) {
                this.push(blk.routeOpts.receiver);
              }
            }, this);
          }, $scope.allReceivers);

          if (!$scope.receivers) {
            var recvs = angular.copy($scope.allReceivers);
            if ($location.search()['receiver']) {
              recvs = angular.copy($location.search()['receiver']);
              // The selected items must always be an array for multi-option selects.
              if (!angular.isArray(recvs)) {
                recvs = [recvs];
              }
            }
            $scope.receivers = recvs;
          }
        },
        function(data) {
          $scope.error = data.data;
        }
      );
    };

    $scope.refresh();
  }
);

angular.module('am.controllers').controller('SilenceCtrl',
  function($scope, $location, Silence) {

    $scope.highlight = $location.search()['hl'] == $scope.sil.id;

    $scope.showDetails = false;
    $scope.showSilenceForm = false;

    $scope.toggleSilenceForm = function() {
      $scope.showSilenceForm = !$scope.showSilenceForm
    }
    $scope.toggleDetails = function() {
      $scope.showDetails = !$scope.showDetails
    }

    var silCopy = angular.copy($scope.sil);

    $scope.delete = function(id) {
      Silence.delete({id: id},
        function(data) {
          $scope.$emit('silence-deleted');
        },
        function(data) {
          $scope.error = data.data;
        });
    };
  }
);

angular.module('am.controllers').controller('SilencesCtrl', function($scope, $location, Silence, silences) {
  $scope.totalItems = silences.data.totalSilences;
  var DEFAULT_PER_PAGE = 25;
  var DEFAULT_PAGE = 1;
  $scope.silences = [];
  $scope.order = "endsAt";

  $scope.showForm = false;

  $scope.toggleForm = function() {
    $scope.showForm = !$scope.showForm;
  }

  // Pagination
  $scope.pageChanged = function() {
    $location.search('page', $scope.currentPage);
  };

  var params = $location.search()
  var page = params['page'];
  if (!page) {
    $location.search('page', DEFAULT_PAGE);
  }
  $scope.currentPage = parseInt($location.search()['page']);

  $scope.itemsPerPage = params['n'];
  if (isNaN(parseInt($scope.itemsPerPage))) {
    $scope.itemsPerPage = DEFAULT_PER_PAGE;
    $location.search('n', $scope.itemsPerPage);
  }

  $scope.setPerPage = function(n) {
    $scope.itemsPerPage = n;
    $location.search('n', $scope.itemsPerPage);
  };

  // Controls the number of pages to display in the pagination list.
  $scope.maxSize = 8;

  // Arbitrary suggested page lengths. The user can override this at any time
  // by entering their own n value in the url.
  $scope.paginationLengths = [5,15,25,50];
  // End Pagination

  $scope.refresh = function() {
    var search = $location.search();
    var params = {};
    if (search['page']) {
      params['offset'] = search['page']-1;
    } else {
      params['offset'] = DEFAULT_PAGE-1;
    }
    $scope.currentPage = params['offset']+1;

    if (search['n']) {
      params['n'] = search['n'];
      $scope.itemsPerPage = params['n'];
    } else {
      params['n'] = $scope.itemsPerPage;
      $scope.setPerPage(params['n']);
    }

    Silence.query(params, function(resp) {
      var data = resp.data;
      $scope.silences = data.silences || [];
      $scope.totalItems = data.totalSilences;
      var now = new Date;

      angular.forEach($scope.silences, function(value) {
        value.endsAt = new Date(value.endsAt);
        value.startsAt = new Date(value.startsAt);
        value.updatedAt = new Date(value.updatedAt);

        value.elapsed = value.endsAt < now;
        value.pending = value.startsAt > now;
        value.active = value.startsAt <= now && value.endsAt > now;
      });
    }, function(data) {
      $scope.error = data.data;
    });
  };

  $scope.$on('silence-created', function(evt) {
    $scope.refresh();
  });
  $scope.$on('silence-deleted', function(evt) {
    $scope.refresh();
  });

  $scope.elapsed = function(elapsed) {
    return function(sil) {
      if (elapsed) {
        return sil.endsAt <= new Date;
      }
      return sil.endsAt > new Date;
    }
  };

  $scope.$watch(function() {
    return $location.search();
  }, function() {
    $scope.refresh();
  }, true);
});

angular.module('am.controllers').controller('SilenceCreateCtrl', function($scope, Silence) {
  $scope.error = null;
  $scope.silence = $scope.silence || {};

  if (!$scope.silence.matchers) {
    $scope.silence.matchers = [{}];
  }

  var origSilence = angular.copy($scope.silence);

  $scope.reset = function() {
    var now = new Date();
    var end = new Date();

    now.setMilliseconds(0);
    end.setMilliseconds(0);
    now.setSeconds(0);
    end.setSeconds(0);

    end.setHours(end.getHours() + 4)

    $scope.silence = angular.copy(origSilence);

    if (!origSilence.startsAt || origSilence.elapsed) {
      $scope.silence.startsAt = now;
    }
    if (!origSilence.endsAt || origSilence.elapsed) {
      $scope.silence.endsAt = end;
    }
  };

  $scope.reset();

  $scope.addMatcher = function() {
    $scope.silence.matchers.push({});
  };

  $scope.delMatcher = function(i) {
    $scope.silence.matchers.splice(i, 1);
  };

  $scope.create = function() {
    var now = new Date;
    // Go through conditions that go against immutability of historic silences.
    var createNew = !angular.equals(origSilence.matchers, $scope.silence.matchers);
    console.log(origSilence, $scope.silence);
    createNew = createNew || $scope.silence.elapsed;
    createNew = createNew || ($scope.silence.active && (origSilence.startsAt == $scope.silence.startsAt || origSilence.endsAt == $scope.silence.endsAt));

    if (createNew) {
      $scope.silence.id = undefined;
    }

    Silence.create($scope.silence, function(data) {
      // If the modifications require creating a new silence,
      // we expire/delete the old one.
      if (createNew && origSilence.id && !$scope.silence.elapsed) {
        Silence.delete({id: origSilence.id},
          function(data) {
            // Only trigger reload after after old silence was deleted.
            $scope.$emit('silence-created');
          },
          function(data) {
            console.warn("deleting silence failed", data);
            $scope.$emit('silence-created');
          });
      } else {
        $scope.$emit('silence-created');
      }
    }, function(data) {
      $scope.error = data.data.error;
    });
  };
});

angular.module('am.services').factory('Status',
  function($resource) {
    return $resource('', {}, {
      'get': {
        method: 'GET',
        url: 'api/v1/status'
      }
    });
  }
);

angular.module('am.controllers').controller('StatusCtrl',
  function($scope, Status) {
    Status.get({},
      function(data) {
        $scope.config = data.data.config;
        $scope.versionInfo = data.data.versionInfo;
        $scope.uptime = data.data.uptime;
      },
      function(data) {
        console.log(data.data);
      })
  }
);

angular.module('am', [
  'ngRoute',
  'ngSanitize',
  'angularMoment',
  'ui.bootstrap',

  'am.controllers',
  'am.services',
  'am.directives'
]);

angular.module('am').config(
  function($routeProvider) {
    $routeProvider.
    when('/alerts', {
      templateUrl: 'app/partials/alerts.html',
      controller: 'AlertsCtrl',
      reloadOnSearch: false
    }).
    when('/silences', {
      templateUrl: 'app/partials/silences.html',
      controller: 'SilencesCtrl',
      resolve: {
        silences: function(Silence) {
          // Required to get the total number of silences before the controller
          // loads. Without this, the user is forced to page 1 of the
          // pagination.
          return Silence.query({'n':0});
        }
      },
      reloadOnSearch: false
    }).
    when('/status', {
      templateUrl: 'app/partials/status.html',
      controller: 'StatusCtrl'
    }).
    otherwise({
      redirectTo: '/alerts'
    });
  }
);
