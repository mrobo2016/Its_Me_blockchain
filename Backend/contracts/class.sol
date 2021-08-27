pragma solidity ^0.5.6;
//pragma experimental ABIEncoderV2;


contract Class {
    enum attendanceStatus {present, late, absent}
    
    struct Attendance {
        string className;
        string classDate;
        uint32 presentPrice;
        uint32 latePrice;
        bool Active; // Enum type
    }

    struct Request {
        address student; 
        string content;
        uint32 reward;
        attendanceStatus status; // Enum type
        bool Active;
    }

    // event - 컨트랙트에서 함수가 실행되는 중간에 이벤트를 발생시켜서, 어떤 변수가 현재 어떤 값인지 로그를 남겨, DApp에서 로그 추적 가능.
    // emit으로 event를 기록 할 수 있음  
    event RequestReceipt(
        uint32 indexed _id,
        address indexed _from,
        string _paymentTxId,
        string _content,
        uint32 _price
    );

    // Global
    address payable owner; // 변수 선언, contstructor에서 initialize 됨 + payable
    string className;

    uint32 numAttendances;
    mapping (uint32 => Attendance) Attendances; // mapping 

    uint32 numRequests;
    mapping (uint32 => Request) requests; // mapping 


    // When a contract is created, its constructor 
    // (a function declared with the constructor keyword) is executed once.
    // A constructor is optional. Only one constructor is allowed, which means overloading is not supported.
    constructor(string memory name, address payable addr) public { // memory keyword -> 동적 할당으로 name을 받는다 
        owner = addr; // payable owner 
        className = name;

        //Attendances[numAttendances] = Attendance("all", 0, 0, MenuStatus.deactivated);
        numAttendances = 0;
        numRequests = 0;
    }


    // public - Public functions are part of the contract interface 
    // and can be either called internally or via messages. For public state variables, an automatic getter function (see below) is generated.
    // view - Functions can be declared view in which case they promise not to modify the state.
    // return 값 -> uint32 dtype
    function getNumAttendances() public view returns (uint32) {
        return numAttendances;
    }

    function getAttendance(uint32 weekId) public view returns (string memory, string memory, uint32, uint32, bool) {
        return (Attendances[weekId].className, Attendances[weekId].classDate, Attendances[weekId].presentPrice, Attendances[weekId].latePrice, Attendances[weekId].Active);
    }

    function addAttendance(string memory classDate, uint32 presentPrice, uint32 latePrice) public{
        require(bytes(classDate).length > 0);
        require(msg.sender == owner);

        Attendances[numAttendances] = Attendance(className, classDate, presentPrice, latePrice, true);
        numAttendances = numAttendances + 1;
        // return attendance number (for confirmation)
    }

    function deactivateAttendance(uint32 weekId) public {
        // require(msg.sender == owner)
        Attendances[weekId].Active = false;
    }


    // Emit
    // Events are emitted using `emit`, followed by
    // the name of the event and the arguments
    // (if any) in parentheses. Any such invocation
    // (even deeply nested) can be detected from
    // the JavaScript API by filtering for `Deposit`.
    function addRequestPresent(string memory paymentTxId, string memory content, uint32 price) public {
        requests[numRequests] = Request(msg.sender, content, price, attendanceStatus.present, true); // msg.sender is a special variable
        emit RequestReceipt(numRequests, msg.sender, paymentTxId, content, price); // emit -> log
        numRequests = numRequests + 1;
    }
    
    function addRequestLate(string memory paymentTxId, string memory content, uint32 price) public {
        requests[numRequests] = Request(msg.sender, content, price, attendanceStatus.late, true); // msg.sender is a special variable
        emit RequestReceipt(numRequests, msg.sender, paymentTxId, content, price); // emit -> log
        numRequests = numRequests + 1;
    }


    function convertAttendanceString(attendanceStatus status) public view returns (string memory){
        if (status == attendanceStatus.present) return "present";
        if (status == attendanceStatus.late) return "late";
        if (status == attendanceStatus.absent) return "absent";
    }
    
    function getRequest(uint32 requestId) public view returns (address, string memory, uint32, string memory, bool) {
        string memory statusString = convertAttendanceString(requests[requestId].status);
        return (requests[requestId].student, requests[requestId].content, requests[requestId].reward, statusString, requests[requestId].Active);
    }

    function approveRequest(uint32 requestId, uint256 reward) public payable {
        require(msg.sender == owner);
        reward = requests[requestId].reward;
        
        if (reward > 0 && requests[requestId].Active) {
            owner.transfer(requests[requestId].reward);  // owner gives student reward for attendance
            requests[requestId].Active = false;
        }
    }

    function denyRequest(uint32 requestId) public {
        requests[requestId].status = attendanceStatus.absent;
        requests[requestId].Active = false;
    }
}
