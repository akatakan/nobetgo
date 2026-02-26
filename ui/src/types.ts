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
  BedCapacity?: number;
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
  IsShiftWorker?: boolean;
  IsActive: boolean;
  Competencies?: string;
  FatigueScore?: number;
  HeroPoint?: number;
  Role?: 'admin' | 'user';
}

export interface AuthResponse {
  token: string;
  role: string;
}

export interface PaginationParams {
  page: number;
  limit: number;
  search?: string;
}

export interface PaginationResult<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface ShiftType extends BaseModel {
  Name: string;
  Description: string;
  StartTime: string;
  EndTime: string;
  Color: string;
  BreakMinutes?: number;
  IsNightShift?: boolean;
  RotationDays?: number;
  DepartmentID?: number;
  Department?: Department;
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

// ===== Time Entry (Puantaj) =====

export interface TimeEntry extends BaseModel {
  employee_id: number;
  employee?: Employee;
  schedule_id?: number;
  schedule?: Schedule;
  clock_in: string;
  clock_out?: string;
  break_minutes: number;
  entry_type: string; // normal, overtime, holiday, weekend
  source: string;     // manual, auto, import
  notes: string;
  status: string;     // pending, approved, rejected
  approved_by?: number;
}

export interface TimeEntryRequest {
  employee_id: number;
  schedule_id?: number;
  clock_in: string;
  clock_out?: string;
  break_minutes?: number;
  entry_type?: string;
  source?: string;
  notes?: string;
}

export interface ClockInRequest {
  employee_id: number;
  notes?: string;
}

export interface ClockOutRequest {
  employee_id: number;
  notes?: string;
}

// ===== Leave (İzin) =====

export interface LeaveType extends BaseModel {
  name: string;
  default_days: number;
  is_paid: boolean;
  requires_approval: boolean;
  color: string;
}

export interface Leave extends BaseModel {
  employee_id: number;
  employee?: Employee;
  leave_type_id: number;
  leave_type?: LeaveType;
  start_date: string;
  end_date: string;
  total_days: number;
  reason: string;
  status: string;
  approved_by?: number;
  approved_at?: string;
}

export interface LeaveRequest {
  employee_id: number;
  leave_type_id: number;
  start_date: string;
  end_date: string;
  reason?: string;
}

export interface LeaveBalance extends BaseModel {
  employee_id: number;
  leave_type_id: number;
  leave_type?: LeaveType;
  year: number;
  total_days: number;
  used_days: number;
  remaining_days: number;
}

// ===== Overtime (Mesai) =====

export interface OvertimeRule extends BaseModel {
  name: string;
  weekly_hour_limit: number;
  daily_hour_limit: number;
  overtime_multiplier: number;
  weekend_multiplier: number;
  holiday_multiplier: number;
  night_shift_extra: number;
  is_active: boolean;
}

export interface PublicHoliday extends BaseModel {
  name: string;
  date: string;
}

export interface OvertimeSummary {
  employee_id: number;
  employee_name: string;
  total_hours: number;
  normal_hours: number;
  overtime_hours: number;
  weekend_hours: number;
  holiday_hours: number;
  night_shift_hours: number;
  working_days: number;
}

// ===== Approval (Onay) =====

export interface AuditLog extends BaseModel {
  entity_type: string;
  entity_id: number;
  action: string;
  field_name?: string;
  old_value?: string;
  new_value?: string;
  performed_by: number;
  ip_address?: string;
}

export interface PendingApprovals {
  time_entries: TimeEntry[];
  leaves: Leave[];
}

// ===== Reports =====

export interface WorkHoursReport {
  month: number;
  year: number;
  employees: EmployeeWorkReport[];
  total_hours: number;
  working_days: number;
}

export interface EmployeeWorkReport {
  employee_id: number;
  employee_name: string;
  department: string;
  total_hours: number;
  working_days: number;
  avg_daily_hours: number;
}

export interface AbsenceReport {
  month: number;
  year: number;
  employees: EmployeeAbsenceReport[];
  total_days: number;
}

export interface EmployeeAbsenceReport {
  employee_id: number;
  employee_name: string;
  leave_type: string;
  total_days: number;
}

export interface EmployeeSummaryReport {
  employee_id: number;
  employee_name: string;
  department: string;
  total_hours: number;
  overtime_hours: number;
  leave_days: number;
  working_days: number;
}

export interface TrendData {
  month: number;
  year: number;
  total_hours: number;
  overtime_hours: number;
  absence_days: number;
  working_days: number;
}

// ===== Form Data (Legacy) =====

export interface ScheduleRequest {
  month: number;
  year: number;
  department_id: number;
  shift_type_ids: number[];
  employee_ids: number[];
  overtime_threshold: number;
  overtime_multiplier: number;
  scheduling_mode?: string;
  beds_per_personnel?: number;
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
  Competencies?: string;
  FatigueScore?: number;
  HeroPoint?: number;
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
  BedCapacity?: number;
}

export interface TitleFormData {
  Name: string;
}

// ===== Cost Analysis =====

export interface CostAnalysisReport {
  month: number;
  year: number;
  employees: EmployeeCostDetail[];
  total_cost: number;
  total_hours: number;
}

export interface EmployeeCostDetail {
  employee_id: number;
  employee_name: string;
  department: string;
  total_hours: number;
  hourly_rate: number;
  total_cost: number;
}

// ===== Notifications =====

export interface Notification extends BaseModel {
  employee_id: number;
  employee?: Employee;
  title: string;
  message: string;
  type: string;
  is_read: boolean;
  action_url?: string;
  related_type?: string;
  related_id?: number;
}
