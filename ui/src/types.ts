// Base model matching GORM's gorm.Model
export interface BaseModel {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
}

export interface Employee extends BaseModel {
  FirstName: string;
  LastName: string;
  Title: string;
  Department: string;
  Email: string;
  Phone: string;
  HourlyRate: number;
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
  Employee: Employee;
  ShiftType: ShiftType;
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
  overtime_threshold: number;
  overtime_multiplier: number;
}

export interface EmployeeFormData {
  FirstName: string;
  LastName: string;
  Title: string;
  Department: string;
  Email: string;
  Phone: string;
  HourlyRate: number;
}

export interface ShiftTypeFormData {
  Name: string;
  StartTime: string;
  EndTime: string;
  Color: string;
  Description: string;
}
