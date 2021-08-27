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

    event RequestReceipt(
        uint32 indexed _id,
        address indexed _from,
        string _paymentTxId,
        string _content,
        uint32 _price
    );

    address payable owner;
    string className;

    uint32 numAttendances;
    mapping (uint32 => Attendance) Attendances; // mapping 

    uint32 numRequests;
    mapping (uint32 => Request) requests; // mapping 

    constructor(string memory name, address payable addr) public { 
        owner = addr; // payable owner 
        className = name;

        //Attendances[numAttendances] = Attendance("all", 0, 0, MenuStatus.deactivated);
        numAttendances = 0;
        numRequests = 0;
    }


    // public - Public functions are part of the contract interface 
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
