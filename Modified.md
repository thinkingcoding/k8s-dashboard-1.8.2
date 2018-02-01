[TOC]
# Modified:
> base from kubernetes dashboard `1.8.2`

## 1.Backend
- [x] add SSO

## 2.Frontend
- [x] add `sessionkeepalive` module
- [x] add **LOGOUT** button for user logout
- [x] change **default namespace** from `default` to `core`

## 3.Usage
### using in yaml
```yaml
...
env:
- name: CLIENT_REDIRECT_URI
  value: https://shclitvm0682.hpeswlab.net:9099
- name: CDF_API_SERVER
  value: https://shclitvm0682.hpeswlab.net:5443
- name: IDM_API_SERVER
  value: https://shclitvm0682.hpeswlab.net:5443
...
```

### end
