import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ModuleConceptComponent } from './module-concept.component';

describe('ModuleConceptComponent', () => {
  let component: ModuleConceptComponent;
  let fixture: ComponentFixture<ModuleConceptComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ModuleConceptComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModuleConceptComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
