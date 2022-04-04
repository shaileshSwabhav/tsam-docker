import { TestBed } from '@angular/core/testing';

import { AddModelComponentService } from './add-model-component.service';

describe('AddModelComponentService', () => {
  let service: AddModelComponentService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AddModelComponentService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
