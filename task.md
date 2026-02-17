# Task Checklist

- [x] **Fix Schedule & Attendance Issues**
  - [x] Remove employee selection limit (verified: no code limit found)
  - [x] Fix overtime calculation (implemented shift-based duration)
  - [x] Allow Floor 0 and negative values in Department (UI & Validation fixed)
  - [x] Verify fixes with E2E test (BrowserSubagent)

- [x] **Manual Attendance for Regular Staff** (New Request)
  - [x] Analyze `AttendanceManager` and `EmployeeManager`
  - [x] Update API to export `employeeApi` and `shiftTypeApi`
  - [x] Update `AttendanceManager.tsx` to include Manual Entry UI
  - [x] **Rename "Nöbet Tipleri" to "Çalışma Tipleri"** <!-- id: 4 -->
- [x] **Employee Manager Improvements** <!-- id: 5 -->
    - [x] Convert Card View to Table View for better scalability
    - [x] Add "IsShiftWorker" flag to Employee model (Backend & Frontend)
    - [x] Filter employees in Schedule Wizard based on "IsShiftWorker"
- [ ] **Data Import Improvements** <!-- id: 6 -->
  - [x] Add `excelize` dependency
  - [x] Implement backend `ImportEmployees` service & handler
  - [x] Update frontend `EmployeeManager` with Import button
  - [x] Update `user_manual.md` with Excel format instructions

- [x] **Documentation**
  - [x] Create `user_manual.md`
  - [x] Update `user_manual.md` with Manual Entry instructions
  - [x] Update `user_manual.md` with Excel Import instructions
  - [x] Create `walkthrough.md` in artifacts
