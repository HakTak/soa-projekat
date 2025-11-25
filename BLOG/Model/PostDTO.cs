public class PostDTO{
        public string Title {get; set;} = null!;

        public string Description {get; set;} = null!;

        public DateTime CreatedAt { get; set; } = DateTime.UtcNow;

        public List<IFormFile> images {get; set;} = null!;  
        
        public int LikeCount {get; set;}
}