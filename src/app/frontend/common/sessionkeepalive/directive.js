/**
 *
 * @param {!angular.$log} $log
 * @param {!angular.$http} $http
 * @param {!./service.SessionService} kdSessionService
 * @return {!angular.Directive}
 * @ngInject
 * */
export default function SessionKeepalive($log, $http, kdSessionService) {
  let lastTime = new Date().getTime();
  return {
    restrict: 'A',
    link: (scope, element, attrs) => {

      let refreshToken = ()=> {
        kdSessionService.setRefreshingState(true);
        console.log('>>>>>>>>>>>>>>>>>>');
        $http.get('refreshtoken').then(
          (/** !angular.$http.Response<Object>*/ response) => {
            let data = angular.toJson(response.data, true);
            console.log(data);
            kdSessionService.setRefreshingState(false);
          },
          (err) => {
            console.log('eeeeeeeeeeeeeeeeeeeeeee');
            kdSessionService.setRefreshingState(false);
          });
      };

      let isTDOA = ()=> {
        return new Date().getTime() - lastTime > 1000 * 60 * 1;
      };

      let checkRefreshSession = () =>{
        if (!kdSessionService.isRefreshing() && isTDOA()) {
          refreshToken();
          lastTime = new Date().getTime();
        }
      };

      element.on('mousedown', () => {
        $log.info('mousedown');
        checkRefreshSession();
      });

      element.on('keydown', () => {
        $log.info('keydown');
        checkRefreshSession();
      });
    }
  };
}


