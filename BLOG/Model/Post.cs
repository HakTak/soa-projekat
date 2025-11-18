
using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace BLOG.Model
{
    public class Post
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string? Id {get; set;}

        [BsonElement("Title")]
        public string Title {get; set;} = null!;

        [BsonElement("Description")]
        public string Description {get; set;} = null!;

        [BsonElement("CreatedAt")]
        public DateTime CreatedAt { get; set; } = DateTime.UtcNow;

        [BsonElement("ImagePaths")]
        public string[]? ImagePaths {get; set;} = null!;

        [BsonElement("LikeCount")]
        public int LikeCount {get; set;}

    }
}