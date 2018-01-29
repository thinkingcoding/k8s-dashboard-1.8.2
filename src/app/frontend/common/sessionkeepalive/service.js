
/**
 *
 * @final
 */
export class SessionService {

  constructor() {
    this.isRefreshingSession = false;
  }

  /**
   * Initializes the service.
   */
  init() {
    this.isRefreshingSession = false;
  }

  /**
   *
   * @export
   * @return {boolean}
   */
  isRefreshing(){
    return this.isRefreshingSession;
  }

  /**
   *
   * @param {boolean} value
   * @export
   */
  setRefreshingState(value){
    this.isRefreshingSession = value;
  }
}
