## WASM plugin for istio using tinygo

cloudevent spec으로 들어오는 http request에서
- cloudevent source key값을 기준으로 
- partitionkey를 추가한다.

spec을 통제할 수 없는 외부에서 cloudevent가 들어오는 걸 kafka에 저장해서 fan out할 때
- partitionkey spec을 활용한다.
- 활성화하면, 같은 partitionkey는 같은 kafka partition으로 들어가기 때문에 source별로 순서가 보장된다.

eventing-kafka의 dispatcher에 istio sidecar를 붙인 뒤 wasm filter를 적용하면
- eventing-kafka 소스코드 수정 없이도 header 정보를 변경할 수 있다.

webAssembly 코드를 만들어야 하는데, 러스트 쓰기 싫어서 tinygo로 구현

### 설치법

mac 기준

```sh
brew tap tinygo-org/tools
brew install tinygo
```


### 테스트

