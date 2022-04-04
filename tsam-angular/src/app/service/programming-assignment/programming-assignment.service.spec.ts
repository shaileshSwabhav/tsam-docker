import { TestBed } from '@angular/core/testing';

import { ProgrammingAssignmentService } from './programming-assignment.service';

describe('ProgrammingAssignmentService', () => {
  let service: ProgrammingAssignmentService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProgrammingAssignmentService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
