// Base model matching GORM's gorm.Model
export interface BaseModel {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
}

export interface Department extends BaseModel {
  Name: string;
  Floor: number;
  Description: string;
}

export interface Title extends BaseModel {
  Name: string;
}

export interface Employee extends BaseModel {
  FirstName: string;
  LastName: string;
  TitleID: number;
  Title: Title;
  DepartmentID: number;
  Department: Department;
  Email: string;
  Phone: string;
  HourlyRate: number;
  IsShiftWorker?: boolean; // Optional for backward compatibility in UI
  IsActive: boolean;
}

export interface ShiftType extends BaseModel {
  Name: string;
  Description: string;
  StartTime: string; // HH:mm
  EndTime: string;   // HH:mm
  Color: string;     // Hex color
}

export interface Schedule extends BaseModel {
  Date: string;
  EmployeeID: number;
  ShiftTypeID: number;
  DepartmentID: number;
  Employee: Employee;
  ShiftType: ShiftType;
  Department: Department;
  IsLocked: boolean;
}

export interface Attendance extends BaseModel {
  ScheduleID: number;
  Schedule: Schedule;
  ActualStartTime: string;
  ActualEndTime: string;
  Notes: string;
  IsOvertime: boolean;
  OvertimeHours: number;
}

export interface ScheduleRequest {
  month: number;
  year: number;
  department_id: number;
  shift_type_ids: number[];
  employee_ids: number[];
  overtime_threshold: number;
  overtime_multiplier: number;
}

export interface EmployeeFormData {
  FirstName: string;
  LastName: string;
  TitleID: number;
  DepartmentID: number;
  Email: string;
  Phone: string;
  HourlyRate: number;
  IsShiftWorker: boolean;
}

export interface ShiftTypeFormData {
  Name: string;
  StartTime: string;
  EndTime: string;
  Color: string;
  Description: string;
}

export interface DepartmentFormData {
  Name: string;
  Floor: number;
  Description: string;
}

export interface TitleFormData {
  Name: string;
}
