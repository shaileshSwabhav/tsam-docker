import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CourseModuleTopicComponent } from './course-module-topic.component';

describe('CourseModuleTopicComponent', () => {
  let component: CourseModuleTopicComponent;
  let fixture: ComponentFixture<CourseModuleTopicComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CourseModuleTopicComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CourseModuleTopicComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
