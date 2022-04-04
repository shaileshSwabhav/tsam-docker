import { TestBed } from '@angular/core/testing';

import { InterviewScheduleService } from './interview-schedule.service';

describe('InterviewScheduleService', () => {
  let service: InterviewScheduleService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(InterviewScheduleService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
