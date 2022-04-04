import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectedTalentListComponent } from './selected-talent-list.component';

describe('SelectedTalentListComponent', () => {
  let component: SelectedTalentListComponent;
  let fixture: ComponentFixture<SelectedTalentListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectedTalentListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectedTalentListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
