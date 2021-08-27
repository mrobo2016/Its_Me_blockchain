pragma solidity ^0.5.6;
pragma experimental ABIEncoderV2;

/*
    - classList는 class 들을 모아둔 저장소이다
    
    - class는 다음과 같은 정보를 담는다
        - professorName -> 교수님 이름
        - className -> 수업명
        - semesterYear -> 개설 학기 + 년도 - 예: "01_2021"
        - classOwner -> 교수가 classOwner이다 [address]
        - classContract -> Class Contract [address]

    - class.sol 과 어떻게 연결되는가?
        - 추정: class.sol contract를 백엔드에서 배포 후, classContract의 address로 전달 받는다. 

*/ 

contract ClassList {
    enum Status {activated, deactivated}
    struct Class {
        string professorName;
        string className;
        string semesterYear;
        address classOwner;
        address classContract;
        Status status;
    }

    address owner;
    Class[] classes;

    constructor() public {
        owner = msg.sender;
    }

    function addClass(string memory professorName, string memory className, string memory semesterYear,
        address classOwner, address classContract) public {
        require(msg.sender == owner);
        // TODO: Input validation will be placed here

        classes.push(Class(professorName, className, semesterYear, classOwner, 
            classContract, Status.activated));
    }

    // gets classes as list 
    function getClasses() public view returns (Class[] memory) {
        return classes;
    }
}