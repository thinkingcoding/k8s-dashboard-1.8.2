/**
 *
 * @param {!angular.$http} $http
 * @param {!./service.SessionService} kdSessionService
 * @return {!angular.Directive}
 * @ngInject
 * */
export default function SessionKeepalive($http, kdSessionService) {
  let lastTime = new Date().getTime();
  return {
    restrict: 'A',
    link: (scope, element, attrs) => {

      let refreshToken = ()=> {
        kdSessionService.startRefresh();
        $http.get('refreshtoken').then(
          () => {
            kdSessionService.endRefresh();
          },
          (err) => {
            kdSessionService.endRefresh();
          });
      };

      let isTDOA = ()=> {
        return new Date().getTime() - lastTime > 1000 * 60 * 5;
      };

      let checkRefreshSession = () =>{
        if (!kdSessionService.isRefreshing() && isTDOA()) {
          refreshToken();
          lastTime = new Date().getTime();
        }
      };

      element.on('mousedown', () => {
        checkRefreshSession();
      });

      element.on('keydown', () => {
        checkRefreshSession();
      });
    }
  };
}


