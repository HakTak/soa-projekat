using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace BLOG.Model
{
    public class BlogLike
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string? Id {get; set;}

        [BsonElement("UserName")]
        public string UserName { get; set; } = null!;
 
        [BsonElement("PostId")]
        public string? BlogId { get; set; } = null!;

        [BsonElement("LikedAt")]
        public DateTime LikedAt { get; set; } = DateTime.UtcNow;

    }
}