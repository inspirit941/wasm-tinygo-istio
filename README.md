## WASM plugin for istio using tinygo

cloudevent spec으로 들어오는 http request에서
- cloudevent source key값을 기준으로 
- partitionkey를 추가한다.

spec을 통제할 수 없는 외부에서 cloudevent가 들어오는 걸 kafka에 저장해서 fan out할 때
- partitionkey spec을 활용한다.
- 활성화하면, 같은 partitionkey는 같은 kafka partition으로 들어가기 때문에 source별로 순서가 보장된다.

eventing-kafka의 dispatcher에 istio sidecar를 붙인 뒤 wasm filter를 적용하면
- eventing-kafka 소스코드 수정 없이도 header 정보를 변경할 수 있다.

특징 / 제약사항

- webAssembly 코드를 만들어야 하는데, 러스트 쓰기 싫어서 tinygo로 구현.
- wasm binary를 따로 관리할 환경이 아니라서, oci image로 관리하기 위해 Dockerfile 사용
  - ko build 같은 다른 옵션도 있지만, Dockerfile을 써야만 하는 제약이 있었다
- eventing-kafka-broker에 istio sidecar붙인 뒤, WasmPlugin extension 적용


### 설치법

mac 기준

```sh
brew tap tinygo-org/tools
brew install tinygo
```


### 테스트

unit test
- `go test` 로 main_test.go 파일 실행.


envoy 적용 테스트

```sh
// wasm binary를 생성한다. plugin.wasm 이라는 binary가 root directory에 생성된다.
make build
 
// envoy container를 띄운다. 빌드 완료한 wasm plugin은 volume으로 envoy container에 포함된다.
docker-compose up

// 로컬에서 요청을 보내고, docker compose의 로그를 확인한다
curl localhost:18000 --header 'ce-specversion: 1.0' --header 'ce-source: //eventsource-1.example.com' 
```


