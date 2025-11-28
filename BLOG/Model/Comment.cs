using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;
using System;

namespace BLOG.Model
{
    public class Comment
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string? Id { get; set; } // jedinstveni ID komentara

        [BsonElement("PostId")]
        public string PostId { get; set; } = null!; // ID blog posta na koji komentar ide

        [BsonElement("AuthorName")]
        public string? AuthorName { get; set; } = null!;

        [BsonElement("Text")]
        public string Text { get; set; } = null!;

        [BsonElement("CreatedAt")]
        public DateTime CreatedAt { get; set; }

        [BsonElement("UpdatedAt")]
        public DateTime UpdatedAt { get; set; }
    }
}
