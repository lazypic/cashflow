![cashflow](./figures/logo.svg)

자본의 흐름, 계획을 도와주는 Backend Commandline Tool

- 용도 : 프리랜서, 그룹, 기업
- 정확한 재무,회계를 도와주는 툴이 아닙니다.
- 4년에 걸쳐 현 시점에서 분기별 매출 및 분기별 지출을 출력하는 간단한 터미널 소프트웨어 입니다.
- 세금을 예측하기 위해서 사용합니다.

### 다운로드
- [Windows 64bit](https://github.com/lazypic/cashflow/releases/download/v0.0.2/cashflow_windows_x86-64.tgz)
- [macOS 64bit](https://github.com/lazypic/cashflow/releases/download/v0.0.2/cashflow_darwin_x86-64.tgz)
- [Linux 64bit](https://github.com/lazypic/cashflow/releases/download/v0.0.2/cashflow_linux_x86-64.tgz)

### 보안
cashflow를 실행하기 위해서는 aws에서 발급된 AccessKey, SecretAccessKey 가 필요합니다.

- Lazypic은 보안정책 설정부분은 내부정책상 프로그래밍 처리하지 않습니다. 각 계정별로 수동처리합니다.

### 사용법
터미널에서 아래처럼 타이핑합니다.
```
$ cashflow
```

출력은 지난해, 이번해, 다음해, 그 다음해의 수입,지출을 분기별로 묶어서 출력합니다.
(단위는 100만원입니다.)
```
+-----+--------+--------+--------+--------+--------+
|     | 2018Q1 | 2018Q2 | 2018Q3 | 2018Q4 | 2018QT |
+-----+--------+--------+--------+--------+--------+
| in  |    0.0 |    0.0 |    0.0 |   14.0 |   14.0 |
| out |    0.0 |    0.0 |    0.0 |    0.0 |    0.0 |
+-----+--------+--------+--------+--------+--------+
+-----+--------+--------+--------+--------+--------+
|     | 2019Q1 | 2019Q2 | 2019Q3 | 2019Q4 | 2019QT |
+-----+--------+--------+--------+--------+--------+
| in  |    6.6 |   14.4 |   18.0 |   24.4 |   63.4 |
| out |    0.0 |    0.0 |    0.0 |    0.0 |    0.0 |
+-----+--------+--------+--------+--------+--------+
+-----+--------+--------+--------+--------+--------+
|     | 2020Q1 | 2020Q2 | 2020Q3 | 2020Q4 | 2020QT |
+-----+--------+--------+--------+--------+--------+
| in  |   15.0 |   18.0 |   30.0 |   31.5 |   74.5 |
| out |    0.0 |    0.0 |    0.0 |    0.0 |    0.0 |
+-----+--------+--------+--------+--------+--------+
+-----+--------+--------+--------+--------+--------+
|     | 2021Q1 | 2021Q2 | 2021Q3 | 2021Q4 | 2021QT |
+-----+--------+--------+--------+--------+--------+
| in  |    4.2 |    0.0 |    0.0 |    0.0 |    4.2 |
| out |    0.0 |    0.0 |    0.0 |    0.0 |    0.0 |
+-----+--------+--------+--------+--------+--------+
```

위 데이터는 실제 데이터가 아닌 임의데이터를 넣었습니다.

> 미수금이 존재한다면, 표가 출력되고 미수금 리스트가 출력됩니다.


### 데이터 입력

기부금 입력시 (기본 Type은 "donation" 입니다.)

```bash
$ cashflow -sender 김한웅 -amount 10000
```

프로젝트 `circle` 계약금 입력시
```bash
$ cashflow -sender 클라이언트 -amount 1000000 -type contract -project circle
```

#### 옵션
- -actualamount : 실제입금금액
- -actualdate : 실제입금일(기본값은 현시간의 RFC3339시간포멧입니다. "2019-05-05T22:30:26+09:00")
- -amount : deposit amount (Required)
- -date : 필수값, 입금일(기본값은 현시간의 RFC3339시간포멧입니다. "2019-05-05T22:30:26+09:00")
- -description : 추가설명
- -help : 도움말 출력
- -profile : AWS db접근을 위해 credentials 파일에 선언된 프로파일명
- -project : 프로젝트코드
- -receivables : 미수금 상태
- -recipient : 받는사람
- -region : db리전. 기본값: 서울리전(ap-northeast-2)
- -sender : 보낸사람
- -table : aws dynamodb 데이터베이스 테이블 이름. 기본값: cashflow
- -type : 입출금 타입. 아래 타입을 사용할 수 있습니다.
	- 수익
		- donation : 기부(기본값)
		- investment : 투자
		- profit : 일시수익
		- contract : 계약금
		- interim, interim1, interim2 : 중도금
		- balance : 잔금
		- addon : 추가금
	- 지출
		- salary : 월급
		- wage : 시급, 주급
		- outsourcing : 외주금
		- other : 활동비
- -unit : 화폐단위. 기본값: `KRW`, [ISO4217](https://en.wikipedia.org/wiki/ISO_4217) 무역 규약에서 허용하는 문자만 입력 가능합니다.

### 백업
dynamoDB는 Full managed DB입니다. 지속 백업기능을 켜서 사용합니다.

### 다국적협업
dynamoDB에서 Global Table을 활성화 시킵니다.

### Software Features
- AWS Serverless를 사용하는 소프트웨어 입니다.
- 1초에 2회 이상 거래가 이용되는 형태로 설계되어있지 않습니다. 하루 최대 86400건의 거래만 저장될 수 있다는 의미 입니다.
- 온디멘드 용량 모드 요금제로 운용됩니다.
	- 대부분의 거래내역은 기가 단위 미만입니다. 백업 서비스를 포함하여 안정적인 서비스가 수년에 걸쳐서 $0~$1의 운용비로 사용할 수 있는 소프트웨어입니다.
- cashflow에서는 데이터의 흐름만 있고, 삭제는 없습니다.
- cashflow는 최종적으로 다음 dynamoDB 들과 연동될 예정입니다.
	- [assetflow](https://github.com/lazypic/assetflow): 회사 에셋(하드웨어, 소프트웨어, 공용계정, 부동산) 비용 관리툴
	- [projectflow](https://github.com/lazypic/projectflow): 프로젝트 비용과 관련된 리스트
	- [userflow](https://github.com/lazypic/userflow): 사용자에게 지출되는 비용과 관련된 리스트
	- [castflow](https://github.com/lazypic/castflow): ipr 관리툴

### 라이센스 : BSD-3-Clause
