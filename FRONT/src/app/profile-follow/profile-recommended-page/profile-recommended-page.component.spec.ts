import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProfileRecommendedPageComponent } from './profile-recommended-page.component';

describe('ProfileRecommendedPageComponent', () => {
  let component: ProfileRecommendedPageComponent;
  let fixture: ComponentFixture<ProfileRecommendedPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ProfileRecommendedPageComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ProfileRecommendedPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
